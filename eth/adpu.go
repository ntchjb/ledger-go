package eth

const (
	ADPU_CLA uint8 = 0xE0

	ADPU_INS_GET_CONFIGURATION         uint8 = 0x06
	ADPU_INS_GET_PUBLIC_KEY            uint8 = 0x02
	ADPU_INS_SIGN_TRANSACTION          uint8 = 0x04
	ADPU_INS_SIGN_PERSONAL_MESSAGE     uint8 = 0x08
	ADPU_INS_SIGN_EIP712               uint8 = 0x0c
	ADPU_INS_ETH2_GET_PUBLIC_KEY       uint8 = 0x0e
	ADPU_INS_ETH2_SET_WITHDRAWAL_INDEX uint8 = 0x10
	ADPU_INS_PRIVACY_OPERATION         uint8 = 0x18

	P1_FIRST_CHUNK uint8 = 0x00
	P1_MORE_CHUNK  uint8 = 0x80

	P1_WITHOUT_CONFIRM uint8 = 0x00
	P1_WITH_CONFIRM    uint8 = 0x01
)
