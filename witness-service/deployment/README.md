# deployment

  deploy.bc contains the bytecode of the contract to deploy for the tube with the following configs:
Safe Address: 0x476c81c27036d05cb5ebfe30ae58c23351a61c4a
    Fee: 20 IOTX
    Min Amount: 1,000 IOTX
    Max Amount: 1,000,000 IOTX

  The safe address could be updated at the deploy time by replacing the address string in the byte code,
or it could be updated after the deployment with "setSafe" api. The gas for deployment is around 1,500,000.

# interact with the tube

  For normal users, "deposit" is the only interface to interact with. The bytecode for this api is
"d0e30db0". With ioctl, assume the contract on the chain is $ca, the following command will transfer 2000
IOTX to the tube:
    ioctl action invoke $ca 2000 -s t1 -l 150000 -p 1 -b d0e30db0
Eventually, user will receive 1980 IOTX token on Ethereum after charing 20 IOTX as the fee.

# FUNCTIONHASHES
    "d0e30db0": "deposit()",
    "67a52793": "depositFee()",
    "f68016b7": "gasLimit()",
    "a0569b57": "getRecords(uint256,uint256)",
    "5f48f393": "maxAmount()",
    "9b2cb5d8": "minAmount()",
    "8da5cb5b": "owner()",
    "8456cb59": "pause()",
    "046f7da2": "resume()",
    "186f0354": "safe()",
    "490ae210": "setDepositFee(uint256)",
    "ee7d72b4": "setGasLimit(uint256)",
    "4fe47f70": "setMaxAmount(uint256)",
    "897b0637": "setMinAmount(uint256)",
    "5db0cb94": "setSafe(address)",
    "f2fde38b": "transferOwnership(address)",
    "2e1a7d4d": "withdraw(uint256)"

