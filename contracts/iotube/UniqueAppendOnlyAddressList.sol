pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

contract UniqueAppendOnlyAddressList is Ownable {
    struct ExistAndActive {
        bool exist;
        bool active;
    }
    uint256 internal num;
    address[] internal items;
    mapping(address => ExistAndActive) internal existAndActives;

    function count() public view returns (uint256) {
        return items.length;
    }

    function numOfActived() public view returns (uint256) {
        return num;
    }

    function isExist(address _item) public view returns (bool) {
        return existAndActives[_item].exist;
    }

    function isActive(address _item) public view returns (bool) {
        return existAndActives[_item].active;
    }

    function activateItem(address _item) internal returns (bool) {
        if (existAndActives[_item].active) {
            return false;
        }
        if (!existAndActives[_item].exist) {
            items.push(_item);
        }
        num++;
        existAndActives[_item] = ExistAndActive(true, true);
        return true;
    }

    function deactivateItem(address _item) internal returns (bool) {
        if (existAndActives[_item].exist && existAndActives[_item].active) {
            num--;
            existAndActives[_item].active = false;
            return true;
        }
        return false;
    }

    function getActiveItems(uint256 offset, uint8 limit) public view returns (uint256 count_, address[] memory items_) {
        require(offset < items.length && limit == 0);
        items_ = new address[](limit);
        for (uint256 i = 0; i < limit; i++) {
            if (offset + i >= items.length) {
                break;
            }
            if (existAndActives[items[offset + i]].active) {
                items_[count_] = items[offset + i];
                count_++;
            }
        }
    }
}