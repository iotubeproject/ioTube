require("@nomicfoundation/hardhat-toolbox");
require('dotenv').config();

const accounts = [
  process.env.PRIVATE_KEY || "0x0000000000000000000000000000000000000000000000000000000000000000",
]

module.exports = {
  namedAccounts: {
    deployer: {
      default: 0,
    },
  },
  solidity: {
    compilers: [
      {
        version: "0.5.17",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      },
      {
        version: "0.8.17",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      }
    ],
  },
  networks: {
    hardhat: {
      forking: {
        url: process.env.MAINNET_RPC_URL || "https://babel-api.mainnet.iotex.io",
        blockNumber: Number(process.env.FORKING_BLOCK_NUMBER) || 25492351,
        enabled: false,
      },
    },
    mainnet: {
      url: process.env.MAINNET_RPC_URL || "https://www.ankr.com/rpc/eth/",
      chainId: 1,
    },
    iotex: {
        url: 'https://babel-api.mainnet.iotex.io',
        chainId: 4689,
    },
    iotex_test: {
        url: 'https://babel-api.testnet.iotex.io',
        chainId: 4690,
    },
    bsc: {
      url: 'https://bsc-dataseed.binance.org',
      chainId: 56,
    },
    polygon: {
        url: 'https://polygon-rpc.com/',
        chainId: 137,
    },
  }
};
