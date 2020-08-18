pragma solidity <0.6 >=0.4.21;

contract Migrations {
  address public owner;
  uint public last_completed_migration;

  modifier restricted() {
    if (msg.sender == owner) _;
  }

  /*@CTK init_migrations
    @post __post.owner == msg.sender
   */
  /* CertiK Smart Labelling, for more details visit: https://certik.org */
  constructor () public {
    owner = msg.sender;
  }

  /*@CTK set_complete
    @pre msg.sender == owner
    @post __post.last_completed_migration == completed
   */
  /* CertiK Smart Labelling, for more details visit: https://certik.org */
  function setCompleted(uint completed) public restricted {
    last_completed_migration = completed;
  }

  function upgrade(address new_address) public restricted {
    Migrations upgraded = Migrations(new_address);
    upgraded.setCompleted(last_completed_migration);
  }
}
