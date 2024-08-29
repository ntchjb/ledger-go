package eip712_test

import (
	"errors"
	"testing"

	"github.com/holiman/uint256"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/stretchr/testify/assert"
)

func Test_Walk(t *testing.T) {
	type fields struct {
		item eip712.StructItem
	}
	type expected struct {
		Type        eip712.Component
		DataCommand eip712.DataCommand
		Path        string
		err         error
	}
	type args struct {
		path string
	}

	err := errors.New("some error")

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected []expected
		err      error
	}{
		{
			name: "Success",
			fields: fields{
				item: eip712.StructItem{
					TypeName: "User",
					Members: []eip712.StructItemMember{
						{
							Name: "id",
							Item: eip712.AtomicItem{
								Item: eip712.NumberData{Num: uint256.NewInt(10), NumBits: 32},
							},
						},
						{
							Name: "fullName",
							Item: eip712.StructItem{
								TypeName: "FullName",
								Members: []eip712.StructItemMember{
									{
										Name: "firstName",
										Item: eip712.AtomicItem{
											Item: eip712.StringData("David"),
										},
									},
									{
										Name: "lastName",
										Item: eip712.AtomicItem{
											Item: eip712.StringData("Ombersburg"),
										},
									},
								},
							},
						},
						{
							Name: "birth",
							Item: eip712.ArrayItem{
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(1990), NumBits: 16},
								},
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(1), NumBits: 8},
								},
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(2), NumBits: 8},
								},
							},
						},
						{
							Name: "occupation",
							Item: eip712.StructItem{
								TypeName: "Occupation",
								Members: []eip712.StructItemMember{
									{
										Name: "job",
										Item: eip712.StructItem{
											TypeName: "Job",
											Members: []eip712.StructItemMember{
												{
													Name: "name",
													Item: eip712.AtomicItem{
														Item: eip712.StringData("developer"),
													},
												},
											},
										},
									},
									{
										Name: "company",
										Item: eip712.AtomicItem{
											Item: eip712.AddressData{
												0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
												0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
											},
										},
									},
									{
										Name: "code",
										Item: eip712.AtomicItem{
											Item: eip712.BytesData{
												Data: []byte{
													0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
												},
												FixedSize: 8,
											},
										},
									},
								},
							},
						},
						{
							Name: "isActive",
							Item: eip712.AtomicItem{
								Item: eip712.BoolData(true),
							},
						},
						{
							Name: "neg",
							Item: eip712.AtomicItem{
								Item: eip712.NumberData{
									Num:     uint256.NewInt(160),
									Signed:  true,
									NumBits: 16,
								},
							},
						},
					},
				},
			},
			args: args{
				path: "",
			},
			expected: []expected{
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("User"),
					},
					Path: "",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x04,
							0x00, 0x00, 0x00, 0x0a,
						},
					},
					Path: "id",
				},
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("FullName"),
					},
					Path: "fullName",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value:     append([]byte{0x00, 0x05}, []byte("David")...),
					},
					Path: "fullName.firstName",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value:     append([]byte{0x00, 0x0A}, []byte("Ombersburg")...),
					},
					Path: "fullName.lastName",
				},
				{
					Type: eip712.DATA_COMPONENT_ARRAY,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ARRAY,
						Value:     []byte{0x03},
					},
					Path: "birth",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x02,
							0x07, 0xC6,
						},
					},
					Path: "birth.[]",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x01,
							0x01,
						},
					},
					Path: "birth.[]",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x01,
							0x02,
						},
					},
					Path: "birth.[]",
				},
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("Occupation"),
					},
					Path: "occupation",
				},
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("Job"),
					},
					Path: "occupation.job",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value:     append([]byte{0x00, 0x09}, []byte("developer")...),
					},
					Path: "occupation.job.name",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x14,
							0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
							0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
						},
					},
					Path: "occupation.company",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x08,
							0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
						},
					},
					Path: "occupation.code",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x01,
							0x01,
						},
					},
					Path: "isActive",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x02,
							0xFF, 0x60,
						},
					},
					Path: "neg",
				},
			},
		},
		{
			name: "Error_StructItemAtRoot_ReturnError",
			fields: fields{
				item: eip712.StructItem{
					TypeName: "User",
					Members: []eip712.StructItemMember{
						{
							Name: "id",
							Item: eip712.AtomicItem{
								Item: eip712.NumberData{Num: uint256.NewInt(10), NumBits: 32},
							},
						},
						{
							Name: "fullName",
							Item: eip712.StructItem{
								TypeName: "FullName",
								Members: []eip712.StructItemMember{
									{
										Name: "firstName",
										Item: eip712.AtomicItem{
											Item: eip712.StringData("David"),
										},
									},
									{
										Name: "lastName",
										Item: eip712.AtomicItem{
											Item: eip712.StringData("Ombersburg"),
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				path: "",
			},
			expected: []expected{
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("User"),
					},
					Path: "",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x04,
							0x00, 0x00, 0x00, 0x0a,
						},
					},
					Path: "id",
				},
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("FullName"),
					},
					Path: "fullName",
					err:  err,
				},
			},
			err: err,
		},
		{
			name: "Error_StructItemAtMember_ReturnError",
			fields: fields{
				item: eip712.StructItem{
					TypeName: "User",
					Members: []eip712.StructItemMember{
						{
							Name: "id",
							Item: eip712.AtomicItem{
								Item: eip712.NumberData{Num: uint256.NewInt(10), NumBits: 32},
							},
						},
						{
							Name: "fullName",
							Item: eip712.StructItem{
								TypeName: "FullName",
								Members: []eip712.StructItemMember{
									{
										Name: "firstName",
										Item: eip712.AtomicItem{
											Item: eip712.StringData("David"),
										},
									},
									{
										Name: "lastName",
										Item: eip712.AtomicItem{
											Item: eip712.StringData("Ombersburg"),
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				path: "",
			},
			expected: []expected{
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("User"),
					},
					Path: "",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x04,
							0x00, 0x00, 0x00, 0x0a,
						},
					},
					Path: "id",
					err:  err,
				},
			},
			err: err,
		},
		{
			name: "Error_ArrayAtMember_ReturnError",
			fields: fields{
				item: eip712.StructItem{
					TypeName: "User",
					Members: []eip712.StructItemMember{
						{
							Name: "id",
							Item: eip712.AtomicItem{
								Item: eip712.NumberData{Num: uint256.NewInt(10), NumBits: 32},
							},
						},
						{
							Name: "birth",
							Item: eip712.ArrayItem{
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(1990), NumBits: 16},
								},
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(1), NumBits: 8},
								},
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(2), NumBits: 8},
								},
							},
						},
					},
				},
			},
			args: args{
				path: "",
			},
			expected: []expected{
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("User"),
					},
					Path: "",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x04,
							0x00, 0x00, 0x00, 0x0a,
						},
					},
					Path: "id",
				},
				{
					Type: eip712.DATA_COMPONENT_ARRAY,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ARRAY,
						Value:     []byte{0x03},
					},
					Path: "birth",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x02,
							0x07, 0xC6,
						},
					},
					Path: "birth.[]",
					err:  err,
				},
			},
			err: err,
		},
		{
			name: "Error_Array_ReturnError",
			fields: fields{
				item: eip712.StructItem{
					TypeName: "User",
					Members: []eip712.StructItemMember{
						{
							Name: "id",
							Item: eip712.AtomicItem{
								Item: eip712.NumberData{Num: uint256.NewInt(10), NumBits: 32},
							},
						},
						{
							Name: "birth",
							Item: eip712.ArrayItem{
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(1990), NumBits: 16},
								},
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(1), NumBits: 8},
								},
								eip712.AtomicItem{
									eip712.NumberData{Num: uint256.NewInt(2), NumBits: 8},
								},
							},
						},
					},
				},
			},
			args: args{
				path: "",
			},
			expected: []expected{
				{
					Type: eip712.DATA_COMPONENT_ROOT,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ROOT,
						Value:     []byte("User"),
					},
					Path: "",
				},
				{
					Type: eip712.DATA_COMPONENT_ATOMIC,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ATOMIC,
						Value: []byte{
							0x00, 0x04,
							0x00, 0x00, 0x00, 0x0a,
						},
					},
					Path: "id",
				},
				{
					Type: eip712.DATA_COMPONENT_ARRAY,
					DataCommand: eip712.DataCommand{
						Component: eip712.DATA_COMPONENT_ARRAY,
						Value:     []byte{0x03},
					},
					Path: "birth",
					err:  err,
				},
			},
			err: err,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			i := 0
			err := test.fields.item.Walk(test.args.path, func(path string, item eip712.Item) error {
				assert.Equal(t, test.expected[i].Type, item.Type())
				assert.Equal(t, test.expected[i].Path, path)
				assert.Equal(t, test.expected[i].DataCommand, item.DataCommand())
				err := test.expected[i].err
				i++

				return err
			})

			assert.ErrorIs(t, err, test.err)
		})
	}
}

func TestDomain_StructItem(t *testing.T) {
	domain := eip712.Domain{
		Name:    "SomeDAppsName",
		Version: "1.0.0",
		ChainID: uint256.NewInt(1),
		VerifyingContract: schema.Address{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		},
		Salt: [32]byte{0x01},
	}

	structItem := domain.StructItem()

	assert.Equal(t, eip712.StructItem{
		TypeName: "EIP712Domain",
		Members: []eip712.StructItemMember{
			{
				Name: "name",
				Item: eip712.AtomicItem{
					Item: eip712.StringData("SomeDAppsName"),
				},
			},
			{
				Name: "version",
				Item: eip712.AtomicItem{
					Item: eip712.StringData("1.0.0"),
				},
			},
			{
				Name: "chainId",
				Item: eip712.AtomicItem{
					Item: eip712.NumberData{
						Num:     uint256.NewInt(1),
						NumBits: 256,
					},
				},
			},
			{
				Name: "verifyingContract",
				Item: eip712.AtomicItem{
					Item: eip712.AddressData{
						0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
						0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
					},
				},
			},
			{
				Name: "salt",
				Item: eip712.AtomicItem{
					Item: eip712.BytesData{
						FixedSize: 32,
						Data: []byte{
							0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
							0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
							0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
							0x00, 0x00,
						},
					},
				},
			},
		},
	}, structItem)
}
