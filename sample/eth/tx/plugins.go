package main

// Derived partially from `https://cdn.live.ledger.com/plugins/ethereum.json`
var plugins = `{
  "0x1111111254eeb25477b68fb85ed929f73a960582": {
    "abi": [
      {
        "inputs": [
          {
            "internalType": "contract IWETH",
            "name": "weth",
            "type": "address"
          }
        ],
        "stateMutability": "nonpayable",
        "type": "constructor"
      },
      { "inputs": [], "name": "AccessDenied", "type": "error" },
      { "inputs": [], "name": "AdvanceNonceFailed", "type": "error" },
      { "inputs": [], "name": "AlreadyFilled", "type": "error" },
      { "inputs": [], "name": "ArbitraryStaticCallFailed", "type": "error" },
      { "inputs": [], "name": "BadPool", "type": "error" },
      { "inputs": [], "name": "BadSignature", "type": "error" },
      { "inputs": [], "name": "ETHTransferFailed", "type": "error" },
      { "inputs": [], "name": "ETHTransferFailed", "type": "error" },
      { "inputs": [], "name": "EmptyPools", "type": "error" },
      { "inputs": [], "name": "EthDepositRejected", "type": "error" },
      { "inputs": [], "name": "GetAmountCallFailed", "type": "error" },
      { "inputs": [], "name": "IncorrectDataLength", "type": "error" },
      { "inputs": [], "name": "InsufficientBalance", "type": "error" },
      { "inputs": [], "name": "InvalidMsgValue", "type": "error" },
      { "inputs": [], "name": "InvalidMsgValue", "type": "error" },
      { "inputs": [], "name": "InvalidatedOrder", "type": "error" },
      { "inputs": [], "name": "MakingAmountExceeded", "type": "error" },
      { "inputs": [], "name": "MakingAmountTooLow", "type": "error" },
      { "inputs": [], "name": "OnlyOneAmountShouldBeZero", "type": "error" },
      { "inputs": [], "name": "OrderExpired", "type": "error" },
      { "inputs": [], "name": "PermitLengthTooLow", "type": "error" },
      { "inputs": [], "name": "PredicateIsNotTrue", "type": "error" },
      { "inputs": [], "name": "PrivateOrder", "type": "error" },
      { "inputs": [], "name": "RFQBadSignature", "type": "error" },
      { "inputs": [], "name": "RFQPrivateOrder", "type": "error" },
      { "inputs": [], "name": "RFQSwapWithZeroAmount", "type": "error" },
      { "inputs": [], "name": "RFQZeroTargetIsForbidden", "type": "error" },
      { "inputs": [], "name": "ReentrancyDetected", "type": "error" },
      { "inputs": [], "name": "RemainingAmountIsZero", "type": "error" },
      { "inputs": [], "name": "ReservesCallFailed", "type": "error" },
      { "inputs": [], "name": "ReturnAmountIsNotEnough", "type": "error" },
      { "inputs": [], "name": "SafePermitBadLength", "type": "error" },
      { "inputs": [], "name": "SafeTransferFailed", "type": "error" },
      { "inputs": [], "name": "SafeTransferFromFailed", "type": "error" },
      {
        "inputs": [
          { "internalType": "bool", "name": "success", "type": "bool" },
          { "internalType": "bytes", "name": "res", "type": "bytes" }
        ],
        "name": "SimulationResults",
        "type": "error"
      },
      { "inputs": [], "name": "SwapAmountTooLarge", "type": "error" },
      { "inputs": [], "name": "SwapWithZeroAmount", "type": "error" },
      { "inputs": [], "name": "TakingAmountExceeded", "type": "error" },
      { "inputs": [], "name": "TakingAmountIncreased", "type": "error" },
      { "inputs": [], "name": "TakingAmountTooHigh", "type": "error" },
      {
        "inputs": [],
        "name": "TransferFromMakerToTakerFailed",
        "type": "error"
      },
      {
        "inputs": [],
        "name": "TransferFromTakerToMakerFailed",
        "type": "error"
      },
      { "inputs": [], "name": "UnknownOrder", "type": "error" },
      { "inputs": [], "name": "WrongAmount", "type": "error" },
      { "inputs": [], "name": "WrongGetter", "type": "error" },
      { "inputs": [], "name": "ZeroAddress", "type": "error" },
      { "inputs": [], "name": "ZeroMinReturn", "type": "error" },
      { "inputs": [], "name": "ZeroReturnAmount", "type": "error" },
      { "inputs": [], "name": "ZeroTargetIsForbidden", "type": "error" },
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": true,
            "internalType": "address",
            "name": "maker",
            "type": "address"
          },
          {
            "indexed": false,
            "internalType": "uint256",
            "name": "newNonce",
            "type": "uint256"
          }
        ],
        "name": "NonceIncreased",
        "type": "event"
      },
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": true,
            "internalType": "address",
            "name": "maker",
            "type": "address"
          },
          {
            "indexed": false,
            "internalType": "bytes32",
            "name": "orderHash",
            "type": "bytes32"
          },
          {
            "indexed": false,
            "internalType": "uint256",
            "name": "remainingRaw",
            "type": "uint256"
          }
        ],
        "name": "OrderCanceled",
        "type": "event"
      },
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": true,
            "internalType": "address",
            "name": "maker",
            "type": "address"
          },
          {
            "indexed": false,
            "internalType": "bytes32",
            "name": "orderHash",
            "type": "bytes32"
          },
          {
            "indexed": false,
            "internalType": "uint256",
            "name": "remaining",
            "type": "uint256"
          }
        ],
        "name": "OrderFilled",
        "type": "event"
      },
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": false,
            "internalType": "bytes32",
            "name": "orderHash",
            "type": "bytes32"
          },
          {
            "indexed": false,
            "internalType": "uint256",
            "name": "makingAmount",
            "type": "uint256"
          }
        ],
        "name": "OrderFilledRFQ",
        "type": "event"
      },
      {
        "anonymous": false,
        "inputs": [
          {
            "indexed": true,
            "internalType": "address",
            "name": "previousOwner",
            "type": "address"
          },
          {
            "indexed": true,
            "internalType": "address",
            "name": "newOwner",
            "type": "address"
          }
        ],
        "name": "OwnershipTransferred",
        "type": "event"
      },
      {
        "inputs": [
          { "internalType": "uint8", "name": "amount", "type": "uint8" }
        ],
        "name": "advanceNonce",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "offsets", "type": "uint256" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "and",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "address", "name": "target", "type": "address" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "arbitraryStaticCall",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "salt", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "receiver",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "offsets",
                "type": "uint256"
              },
              {
                "internalType": "bytes",
                "name": "interactions",
                "type": "bytes"
              }
            ],
            "internalType": "struct OrderLib.Order",
            "name": "order",
            "type": "tuple"
          }
        ],
        "name": "cancelOrder",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "orderRemaining",
            "type": "uint256"
          },
          { "internalType": "bytes32", "name": "orderHash", "type": "bytes32" }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "orderInfo", "type": "uint256" }
        ],
        "name": "cancelOrderRFQ",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "orderInfo", "type": "uint256" },
          {
            "internalType": "uint256",
            "name": "additionalMask",
            "type": "uint256"
          }
        ],
        "name": "cancelOrderRFQ",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "salt", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "receiver",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "offsets",
                "type": "uint256"
              },
              {
                "internalType": "bytes",
                "name": "interactions",
                "type": "bytes"
              }
            ],
            "internalType": "struct OrderLib.Order",
            "name": "order",
            "type": "tuple"
          }
        ],
        "name": "checkPredicate",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "contract IClipperExchangeInterface",
            "name": "clipperExchange",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "dstToken",
            "type": "address"
          },
          {
            "internalType": "uint256",
            "name": "inputAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "outputAmount",
            "type": "uint256"
          },
          { "internalType": "uint256", "name": "goodUntil", "type": "uint256" },
          { "internalType": "bytes32", "name": "r", "type": "bytes32" },
          { "internalType": "bytes32", "name": "vs", "type": "bytes32" }
        ],
        "name": "clipperSwap",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "contract IClipperExchangeInterface",
            "name": "clipperExchange",
            "type": "address"
          },
          {
            "internalType": "address payable",
            "name": "recipient",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "dstToken",
            "type": "address"
          },
          {
            "internalType": "uint256",
            "name": "inputAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "outputAmount",
            "type": "uint256"
          },
          { "internalType": "uint256", "name": "goodUntil", "type": "uint256" },
          { "internalType": "bytes32", "name": "r", "type": "bytes32" },
          { "internalType": "bytes32", "name": "vs", "type": "bytes32" }
        ],
        "name": "clipperSwapTo",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "contract IClipperExchangeInterface",
            "name": "clipperExchange",
            "type": "address"
          },
          {
            "internalType": "address payable",
            "name": "recipient",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "dstToken",
            "type": "address"
          },
          {
            "internalType": "uint256",
            "name": "inputAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "outputAmount",
            "type": "uint256"
          },
          { "internalType": "uint256", "name": "goodUntil", "type": "uint256" },
          { "internalType": "bytes32", "name": "r", "type": "bytes32" },
          { "internalType": "bytes32", "name": "vs", "type": "bytes32" },
          { "internalType": "bytes", "name": "permit", "type": "bytes" }
        ],
        "name": "clipperSwapToWithPermit",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [],
        "name": "destroy",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "value", "type": "uint256" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "eq",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "salt", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "receiver",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "offsets",
                "type": "uint256"
              },
              {
                "internalType": "bytes",
                "name": "interactions",
                "type": "bytes"
              }
            ],
            "internalType": "struct OrderLib.Order",
            "name": "order",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "signature", "type": "bytes" },
          { "internalType": "bytes", "name": "interaction", "type": "bytes" },
          {
            "internalType": "uint256",
            "name": "makingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "takingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "skipPermitAndThresholdAmount",
            "type": "uint256"
          }
        ],
        "name": "fillOrder",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "bytes32", "name": "", "type": "bytes32" }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "info", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              }
            ],
            "internalType": "struct OrderRFQLib.OrderRFQ",
            "name": "order",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "signature", "type": "bytes" },
          {
            "internalType": "uint256",
            "name": "flagsAndAmount",
            "type": "uint256"
          }
        ],
        "name": "fillOrderRFQ",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "bytes32", "name": "", "type": "bytes32" }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "info", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              }
            ],
            "internalType": "struct OrderRFQLib.OrderRFQ",
            "name": "order",
            "type": "tuple"
          },
          { "internalType": "bytes32", "name": "r", "type": "bytes32" },
          { "internalType": "bytes32", "name": "vs", "type": "bytes32" },
          {
            "internalType": "uint256",
            "name": "flagsAndAmount",
            "type": "uint256"
          }
        ],
        "name": "fillOrderRFQCompact",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "filledMakingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "filledTakingAmount",
            "type": "uint256"
          },
          { "internalType": "bytes32", "name": "orderHash", "type": "bytes32" }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "info", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              }
            ],
            "internalType": "struct OrderRFQLib.OrderRFQ",
            "name": "order",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "signature", "type": "bytes" },
          {
            "internalType": "uint256",
            "name": "flagsAndAmount",
            "type": "uint256"
          },
          { "internalType": "address", "name": "target", "type": "address" }
        ],
        "name": "fillOrderRFQTo",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "filledMakingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "filledTakingAmount",
            "type": "uint256"
          },
          { "internalType": "bytes32", "name": "orderHash", "type": "bytes32" }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "info", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              }
            ],
            "internalType": "struct OrderRFQLib.OrderRFQ",
            "name": "order",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "signature", "type": "bytes" },
          {
            "internalType": "uint256",
            "name": "flagsAndAmount",
            "type": "uint256"
          },
          { "internalType": "address", "name": "target", "type": "address" },
          { "internalType": "bytes", "name": "permit", "type": "bytes" }
        ],
        "name": "fillOrderRFQToWithPermit",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "bytes32", "name": "", "type": "bytes32" }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "salt", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "receiver",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "offsets",
                "type": "uint256"
              },
              {
                "internalType": "bytes",
                "name": "interactions",
                "type": "bytes"
              }
            ],
            "internalType": "struct OrderLib.Order",
            "name": "order_",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "signature", "type": "bytes" },
          { "internalType": "bytes", "name": "interaction", "type": "bytes" },
          {
            "internalType": "uint256",
            "name": "makingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "takingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "skipPermitAndThresholdAmount",
            "type": "uint256"
          },
          { "internalType": "address", "name": "target", "type": "address" }
        ],
        "name": "fillOrderTo",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "actualMakingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "actualTakingAmount",
            "type": "uint256"
          },
          { "internalType": "bytes32", "name": "orderHash", "type": "bytes32" }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "salt", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "receiver",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "offsets",
                "type": "uint256"
              },
              {
                "internalType": "bytes",
                "name": "interactions",
                "type": "bytes"
              }
            ],
            "internalType": "struct OrderLib.Order",
            "name": "order",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "signature", "type": "bytes" },
          { "internalType": "bytes", "name": "interaction", "type": "bytes" },
          {
            "internalType": "uint256",
            "name": "makingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "takingAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "skipPermitAndThresholdAmount",
            "type": "uint256"
          },
          { "internalType": "address", "name": "target", "type": "address" },
          { "internalType": "bytes", "name": "permit", "type": "bytes" }
        ],
        "name": "fillOrderToWithPermit",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "uint256", "name": "", "type": "uint256" },
          { "internalType": "bytes32", "name": "", "type": "bytes32" }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "value", "type": "uint256" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "gt",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "components": [
              { "internalType": "uint256", "name": "salt", "type": "uint256" },
              {
                "internalType": "address",
                "name": "makerAsset",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "takerAsset",
                "type": "address"
              },
              { "internalType": "address", "name": "maker", "type": "address" },
              {
                "internalType": "address",
                "name": "receiver",
                "type": "address"
              },
              {
                "internalType": "address",
                "name": "allowedSender",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "makingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "takingAmount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "offsets",
                "type": "uint256"
              },
              {
                "internalType": "bytes",
                "name": "interactions",
                "type": "bytes"
              }
            ],
            "internalType": "struct OrderLib.Order",
            "name": "order",
            "type": "tuple"
          }
        ],
        "name": "hashOrder",
        "outputs": [
          { "internalType": "bytes32", "name": "", "type": "bytes32" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [],
        "name": "increaseNonce",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "address", "name": "maker", "type": "address" },
          { "internalType": "uint256", "name": "slot", "type": "uint256" }
        ],
        "name": "invalidatorForOrderRFQ",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "value", "type": "uint256" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "lt",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "address", "name": "", "type": "address" }
        ],
        "name": "nonce",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "address",
            "name": "makerAddress",
            "type": "address"
          },
          { "internalType": "uint256", "name": "makerNonce", "type": "uint256" }
        ],
        "name": "nonceEquals",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "offsets", "type": "uint256" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "or",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [],
        "name": "owner",
        "outputs": [
          { "internalType": "address", "name": "", "type": "address" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "bytes32", "name": "orderHash", "type": "bytes32" }
        ],
        "name": "remaining",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "bytes32", "name": "orderHash", "type": "bytes32" }
        ],
        "name": "remainingRaw",
        "outputs": [
          { "internalType": "uint256", "name": "", "type": "uint256" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "bytes32[]",
            "name": "orderHashes",
            "type": "bytes32[]"
          }
        ],
        "name": "remainingsRaw",
        "outputs": [
          { "internalType": "uint256[]", "name": "", "type": "uint256[]" }
        ],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [],
        "name": "renounceOwnership",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "contract IERC20",
            "name": "token",
            "type": "address"
          },
          { "internalType": "uint256", "name": "amount", "type": "uint256" }
        ],
        "name": "rescueFunds",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "address", "name": "target", "type": "address" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "simulate",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "contract IAggregationExecutor",
            "name": "executor",
            "type": "address"
          },
          {
            "components": [
              {
                "internalType": "contract IERC20",
                "name": "srcToken",
                "type": "address"
              },
              {
                "internalType": "contract IERC20",
                "name": "dstToken",
                "type": "address"
              },
              {
                "internalType": "address payable",
                "name": "srcReceiver",
                "type": "address"
              },
              {
                "internalType": "address payable",
                "name": "dstReceiver",
                "type": "address"
              },
              {
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
              },
              {
                "internalType": "uint256",
                "name": "minReturnAmount",
                "type": "uint256"
              },
              { "internalType": "uint256", "name": "flags", "type": "uint256" }
            ],
            "internalType": "struct GenericRouter.SwapDescription",
            "name": "desc",
            "type": "tuple"
          },
          { "internalType": "bytes", "name": "permit", "type": "bytes" },
          { "internalType": "bytes", "name": "data", "type": "bytes" }
        ],
        "name": "swap",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "spentAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "time", "type": "uint256" }
        ],
        "name": "timestampBelow",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "uint256",
            "name": "timeNonceAccount",
            "type": "uint256"
          }
        ],
        "name": "timestampBelowAndNonceEquals",
        "outputs": [{ "internalType": "bool", "name": "", "type": "bool" }],
        "stateMutability": "view",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "address", "name": "newOwner", "type": "address" }
        ],
        "name": "transferOwnership",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          { "internalType": "uint256", "name": "amount", "type": "uint256" },
          { "internalType": "uint256", "name": "minReturn", "type": "uint256" },
          { "internalType": "uint256[]", "name": "pools", "type": "uint256[]" }
        ],
        "name": "uniswapV3Swap",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "int256",
            "name": "amount0Delta",
            "type": "int256"
          },
          {
            "internalType": "int256",
            "name": "amount1Delta",
            "type": "int256"
          },
          { "internalType": "bytes", "name": "", "type": "bytes" }
        ],
        "name": "uniswapV3SwapCallback",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "address payable",
            "name": "recipient",
            "type": "address"
          },
          { "internalType": "uint256", "name": "amount", "type": "uint256" },
          { "internalType": "uint256", "name": "minReturn", "type": "uint256" },
          { "internalType": "uint256[]", "name": "pools", "type": "uint256[]" }
        ],
        "name": "uniswapV3SwapTo",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "address payable",
            "name": "recipient",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          { "internalType": "uint256", "name": "amount", "type": "uint256" },
          { "internalType": "uint256", "name": "minReturn", "type": "uint256" },
          { "internalType": "uint256[]", "name": "pools", "type": "uint256[]" },
          { "internalType": "bytes", "name": "permit", "type": "bytes" }
        ],
        "name": "uniswapV3SwapToWithPermit",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          { "internalType": "uint256", "name": "amount", "type": "uint256" },
          { "internalType": "uint256", "name": "minReturn", "type": "uint256" },
          { "internalType": "uint256[]", "name": "pools", "type": "uint256[]" }
        ],
        "name": "unoswap",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "address payable",
            "name": "recipient",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          { "internalType": "uint256", "name": "amount", "type": "uint256" },
          { "internalType": "uint256", "name": "minReturn", "type": "uint256" },
          { "internalType": "uint256[]", "name": "pools", "type": "uint256[]" }
        ],
        "name": "unoswapTo",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "payable",
        "type": "function"
      },
      {
        "inputs": [
          {
            "internalType": "address payable",
            "name": "recipient",
            "type": "address"
          },
          {
            "internalType": "contract IERC20",
            "name": "srcToken",
            "type": "address"
          },
          { "internalType": "uint256", "name": "amount", "type": "uint256" },
          { "internalType": "uint256", "name": "minReturn", "type": "uint256" },
          { "internalType": "uint256[]", "name": "pools", "type": "uint256[]" },
          { "internalType": "bytes", "name": "permit", "type": "bytes" }
        ],
        "name": "unoswapToWithPermit",
        "outputs": [
          {
            "internalType": "uint256",
            "name": "returnAmount",
            "type": "uint256"
          }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
      },
      { "stateMutability": "payable", "type": "receive" }
    ],
    "0x0502b1c5": {
      "erc20OfInterest": ["srcToken"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a9605820502b1c5",
      "signature": "304402204561f90c5dbb09e2aa0a748a4810a3ebf263478cfb554d7587169808bf09af4702202bfc682db8cdd19b904f17eb40a90a16c9bfacf87535347c5c0804238f0e8037"
    },
    "0x12aa3caf": {
      "erc20OfInterest": ["desc.srcToken", "desc.dstToken"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a96058212aa3caf",
      "signature": "30450221009bf7192ed1276263000f619b6133c98a393bff309ac8901b5593849fbf276b2702202de029f07bd0573737b368a80d592daa331ad982b0b9162ecc301fb017846e8c"
    },
    "0x3c15fd91": {
      "erc20OfInterest": ["srcToken"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a9605823c15fd91",
      "signature": "3045022100b2331a720af1009502ffcda388b8dce452cc89ae0c78a5b51eb5f57b83df7ba00220643cf85b11205f616fab2f1cfbe423e5bbd67dab5e5731c274ba334af4e1418e"
    },
    "0x3eca9c0a": {
      "erc20OfInterest": ["order.makerAsset", "order.takerAsset"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a9605823eca9c0a",
      "signature": "3045022100a13d07c910677e5825daa40cb2189f64dd77b1e04e23898cbe89812c564e314f022003c64847141c76c494ea1bb773967bd4b2a4b8978d7bfc917aebaba3b1bd33dd"
    },
    "0x70ccbd31": {
      "erc20OfInterest": ["order.makerAsset", "order.takerAsset"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a96058270ccbd31",
      "signature": "3045022100ed214278eb77d5bb5b6b02f85c9685e809c4dd508fce0675e36281c84f57baa1022070ccdf82ddf4de764dbf9d61a6849753e2f378d6b66f5ce6c41430e9b8d13963"
    },
    "0x84bd6d29": {
      "erc20OfInterest": ["srcToken", "dstToken"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a96058284bd6d29",
      "signature": "304402203bbbf743963cb776a298ddbbfeecad7d7180f614f36210dee1d03caab66508c7022003cdf7ee2656b05c663c5775fb04436075f69800bb423b22e27e7d4896471864"
    },
    "0xc805a666": {
      "erc20OfInterest": ["srcToken", "dstToken"],
      "plugin": "1inch",
      "serialized_data": "0531696e63681111111254eeb25477b68fb85ed929f73a960582c805a666",
      "signature": "30440220328be01ca781cae549ed0320511d7b934024fa70797ac1273c8be89e1feb8b6b02201e3c79dfb6b38497fc2df6fbcd4804c7ec814ef63e035144754d5c19c5e6cb44"
    }
  }
}`
