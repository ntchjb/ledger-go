package rlp_test

import (
	"testing"

	"github.com/ntchjb/ledger-go/eth/rlp"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		encoded []byte
		res     rlp.Item
		n       int
		err     error
	}{
		{
			name: "Success_ShortString",
			encoded: []byte{
				0x83, 'd', 'o', 'g',
			},
			res: rlp.Item{
				Data: []byte("dog"),
			},
			n:   4,
			err: nil,
		},
		{
			name: "Success_ShortList",
			encoded: []byte{
				0xc8,
				0x83, 'c', 'a', 't',
				0x83, 'd', 'o', 'g',
			},
			res: rlp.Item{
				List: []rlp.Item{
					{
						Data: []byte("cat"),
					},
					{
						Data: []byte("dog"),
					},
				},
			},
			n:   9,
			err: nil,
		},
		{
			name: "Success_EmptyString",
			encoded: []byte{
				0x80,
			},
			res: rlp.Item{
				Data: []byte{},
			},
			n:   1,
			err: nil,
		},
		{
			name: "Success_EmptyList",
			encoded: []byte{
				0xC0,
			},
			res: rlp.Item{
				List: []rlp.Item{},
			},
			n:   1,
			err: nil,
		},
		{
			name: "Success_Byte0",
			encoded: []byte{
				0x00,
			},
			res: rlp.Item{
				Data: []byte{0x00},
			},
			n:   1,
			err: nil,
		},
		{
			name: "Success_CharacterA",
			encoded: []byte{
				'a',
			},
			res: rlp.Item{
				Data: []byte{'a'},
			},
			n:   1,
			err: nil,
		},
		{
			name: "Success_NestedEmptyArray",
			encoded: []byte{
				0xC7, 0xC0, 0xC1, 0xC0, 0xC3, 0xC0, 0xC1, 0xC0,
			},
			res: rlp.Item{
				List: []rlp.Item{
					{
						List: []rlp.Item{},
					},
					{
						List: []rlp.Item{
							{
								List: []rlp.Item{},
							},
						},
					},
					{
						List: []rlp.Item{
							{
								List: []rlp.Item{},
							},
							{
								List: []rlp.Item{
									{
										List: []rlp.Item{},
									},
								},
							},
						},
					},
				},
			},
			n:   8,
			err: nil,
		},
		{
			name: "Success_Complex",
			encoded: []byte{
				0xF8,
				0x44,
				0xC0,
				0xB8, 0x38,
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u',
				'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o',
				'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'g', ' ', 'e', 'l', 'i', 't',
				0xC1, 0xC0,
				0xC6, 'a', 0xC4, 0x83, 'd', 'o', 'g',
			},
			res: rlp.Item{
				List: []rlp.Item{
					{
						List: []rlp.Item{},
					},
					{
						Data: []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
					},
					{
						List: []rlp.Item{
							{
								List: []rlp.Item{},
							},
						},
					},
					{
						List: []rlp.Item{
							{
								Data: []byte{'a'},
							},
							{
								List: []rlp.Item{
									{
										Data: []byte("dog"),
									},
								},
							},
						},
					},
				},
			},
			n:   70,
			err: nil,
		},
		{
			name: "Success_LongString",
			encoded: []byte{
				0xB8, 0x38,
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u',
				'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o',
				'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'g', ' ', 'e', 'l', 'i', 't',
			},
			res: rlp.Item{
				Data: []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
			},
			n:   58,
			err: nil,
		},
		{
			name: "Success_LongString2",
			encoded: []byte{
				0xB9, 0x01, 0x02,
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'L', 'o', 'r', 'e', 'm', ' ', 'i', 'p', 's', 'u', 'm', ' ', 'd', 'o', 'l', 'o', 'r', ' ', 's', 'i',
				't', ' ', 'a', 'm', 'e', 't', ',', ' ', 'c', 'o', 'n', 's', 'e', 'c', 't', 'e', 't', 'u', 'r', ' ',
				'a', 'd', 'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
				'i', 'p', 'i', 's', 'i', 'c', 'i', 'n',
			},
			res: rlp.Item{
				Data: []byte(
					"Lorem ipsum dolor sit amet, consectetur adipisicin" +
						"Lorem ipsum dolor sit amet, consectetur adipisicin" +
						"Lorem ipsum dolor sit amet, consectetur adipisicin" +
						"Lorem ipsum dolor sit amet, consectetur adipisicin" +
						"Lorem ipsum dolor sit amet, consectetur adipisicin" +
						"ipisicin",
				),
			},
			n:   261,
			err: nil,
		},
		{
			name: "Success_LongList",
			encoded: []byte{
				0xF8, 0x3c,
				0x84, 'd', 'o', 'g', '0',
				0x84, 'd', 'o', 'g', '1',
				0x84, 'd', 'o', 'g', '2',
				0x84, 'd', 'o', 'g', '3',
				0x84, 'd', 'o', 'g', '4',
				0x84, 'd', 'o', 'g', '5',
				0x84, 'd', 'o', 'g', '6',
				0x84, 'd', 'o', 'g', '7',
				0x84, 'd', 'o', 'g', '8',
				0x84, 'd', 'o', 'g', '9',
				0x84, 'd', 'o', 'o', '0',
				0x84, 'd', 'o', 'o', '1',
			},
			res: rlp.Item{
				List: []rlp.Item{
					{Data: []byte("dog0")},
					{Data: []byte("dog1")},
					{Data: []byte("dog2")},
					{Data: []byte("dog3")},
					{Data: []byte("dog4")},
					{Data: []byte("dog5")},
					{Data: []byte("dog6")},
					{Data: []byte("dog7")},
					{Data: []byte("dog8")},
					{Data: []byte("dog9")},
					{Data: []byte("doo0")},
					{Data: []byte("doo1")},
				},
			},
			n:   62,
			err: nil,
		},
		{
			name: "Success_LongList2",
			encoded: []byte{
				0xF9, 0x01, 0x01,
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '0',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '1',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '2',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '3',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '4',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '5',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '6',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '7',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '8',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '0', '9',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '1', '0',
				0x93, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm', 'e', '1', '1',
				0x90, 'L', 'o', 'r', 'e', 'm', 'I', 'p', 's', 'u', 'm', 'A', 'w', 'e', 's', 'o', 'm',
			},
			res: rlp.Item{
				List: []rlp.Item{
					{Data: []byte("LoremIpsumAwesome00")},
					{Data: []byte("LoremIpsumAwesome01")},
					{Data: []byte("LoremIpsumAwesome02")},
					{Data: []byte("LoremIpsumAwesome03")},
					{Data: []byte("LoremIpsumAwesome04")},
					{Data: []byte("LoremIpsumAwesome05")},
					{Data: []byte("LoremIpsumAwesome06")},
					{Data: []byte("LoremIpsumAwesome07")},
					{Data: []byte("LoremIpsumAwesome08")},
					{Data: []byte("LoremIpsumAwesome09")},
					{Data: []byte("LoremIpsumAwesome10")},
					{Data: []byte("LoremIpsumAwesome11")},
					{Data: []byte("LoremIpsumAwesom")},
				},
			},
			n:   260,
			err: nil,
		},
		{
			name:    "Error_InputNil",
			encoded: []byte{},
			res:     rlp.Item{},
			n:       0,
			err:     rlp.ErrEmptyInput,
		},
		{
			name: "Error_InvalidFormat",
			encoded: []byte{
				0x83, 'd', 'o',
			},
			res: rlp.Item{},
			n:   0,
			err: rlp.ErrInvalidRLPFormat,
		},
		{
			name: "Error_InvalidFormatNested",

			encoded: []byte{
				0xC7, 0xC0, 0xC1, 0xC0, 0xC3, 0x83, 'd', 'o',
			},
			res: rlp.Item{},
			n:   0,
			err: rlp.ErrInvalidRLPFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, n, gotErr := rlp.Decode(test.encoded)

			assert.Equal(t, test.res, actual)
			assert.Equal(t, test.n, n)
			assert.Equal(t, test.n, test.res.Len())
			assert.ErrorIs(t, gotErr, test.err)
		})
	}

}
