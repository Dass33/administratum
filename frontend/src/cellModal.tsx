import React, { useEffect, useRef, useState } from 'react';
import { useApp, EnumColTypes, Column, ColumnData, DEFAULT_UUID, Sheet, NullString, EnumSheetTypes, ColTypes, Domain } from './AppContext';
import Dropdown from "./dropdown";
import { DropdownOption } from "./dropdown";
import cross from "./assets/cross.svg";

interface ArrayItem {
    value: string;
    isValid: boolean;
}

const CellModal = () => {
    const {
        setCellModal, cellModal,
        columns, setColumns,
        accessToken,
        currSheet, setCurrSheet,
        currBranch,
    } = useApp();
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    const numberInputRef = useRef<HTMLInputElement>(null);

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const rowIdx = cellModal
        ? cellModal[0]
        : -1;

    const currCol = cellModal
        ? cellModal[1]
        : null;

    const initCell = currCol?.data.find(item => item.idx == rowIdx)
    const initCellVal = initCell
        ? initCell.value.String
        : ""

    const isConfig = currSheet?.type == EnumSheetTypes.MAP;

    const [cellVal, setCellVal] = useState(initCellVal);
    const [boolVal, setBoolVal] = useState(initCellVal == 'true');
    const [arrayItems, setArrayItems] = useState<ArrayItem[]>([]);
    const [arrayType, setArrayType] = useState<EnumColTypes.TEXT | EnumColTypes.NUMBER>(EnumColTypes.TEXT);
    const [enumVal, setEnumVal] = useState(initCellVal);

    const [colType, setColType] = useState<string>(() => {
        if (isConfig) {
            if (initCell?.type.Valid) return initCell.type.String;
            return EnumColTypes.TEXT;
        }
        if (cellModal) return cellModal[1].type
        return EnumColTypes.TEXT;
    })

    const optionsBranches: DropdownOption[] = [
        { value: EnumColTypes.TEXT, label: EnumColTypes.TEXT },
        { value: EnumColTypes.NUMBER, label: EnumColTypes.NUMBER }
    ];

    const isEnumType = (type: string): boolean => {
        const baseTypes = [EnumColTypes.TEXT, EnumColTypes.NUMBER, EnumColTypes.BOOL, EnumColTypes.ARRAY];
        return !baseTypes.includes(type as EnumColTypes);
    };

    const getEnumValues = (enumName: string): string[] => {
        const enumItem = currBranch?.enums?.find(e => e.name === enumName);
        return enumItem?.vals || [];
    };

    const enumOptions = isEnumType(colType) ? getEnumValues(colType).map(val => ({ 
        value: val, 
        label: val 
    })) : [];
    

    useEffect(() => {
        if (colType === EnumColTypes.ARRAY && initCellVal) {
            try {
                const parsed = JSON.parse(initCellVal);
                if (Array.isArray(parsed)) {
                    const firstItem = parsed[0];
                    const detectedType = typeof firstItem === 'number' ? EnumColTypes.NUMBER : EnumColTypes.TEXT;
                    setArrayType(detectedType);

                    const items: ArrayItem[] = parsed.map((item, _) => ({
                        value: String(item),
                        isValid: true
                    }));
                    setArrayItems(items);
                }
            } catch (e) {
                setArrayItems([]);
            }
        }
    }, [colType, initCellVal]);

    const validateArrayItem = (value: string, type: EnumColTypes.TEXT | EnumColTypes.NUMBER): boolean => {
        if (type === EnumColTypes.NUMBER) {
            return !isNaN(Number(value)) && value.trim() !== '';
        }
        return true;
    };

    const updateArrayItem = (index: number, field: keyof ArrayItem, value: any) => {
        setArrayItems(prev => {
            const newItems = [...prev];
            newItems[index] = { ...newItems[index], [field]: value };

            if (field === 'value') {
                newItems[index].isValid = validateArrayItem(value, arrayType);
            }

            return newItems;
        });
    };

    const addArrayItem = () => {
        setArrayItems(prev => [...prev, {
            value: '',
            isValid: true
        }]);
    };

    const removeArrayItem = (index: number) => {
        setArrayItems(prev => prev.filter((_, i) => i !== index));
    };

    const updateCell = (newVal: NullString, rowIndex: number, col: Column) => {
        if (!currSheet) return;
        let item_id = DEFAULT_UUID;
        let itemFound = false;
        const newCol = { ...col };

        newCol.data = col.data.map(item => {
            if (item.idx === rowIndex) {
                itemFound = true;
                item_id = item.id;
                return {
                    id: item_id,
                    idx: rowIndex,
                    value: newVal,
                    type: { String: colType, Valid: isConfig },
                };
            }
            return item;
        });

        if (itemFound) {
            const updatedData: ColumnData = {
                id: item_id,
                idx: rowIndex,
                value: newVal,
                type: { String: colType, Valid: isConfig },
            };
            putAdjustedColumnData(updatedData, accessToken ?? "");
        } else {
            const newColData: ColumnData = {
                id: item_id,
                idx: rowIndex,
                value: newVal,
                type: { String: colType, Valid: isConfig },
            };
            newCol.data.push(newColData);
            postNewColumnData(col, newColData, currSheet, accessToken ?? "");

            const newSheet = currSheet;
            if (newSheet && newSheet.row_count <= rowIndex) newSheet.row_count++;
            setCurrSheet(newSheet);
        }

        const newColumns = columns.map(item => {
            if (item.name === col.name) return newCol;
            return item;
        });
        setColumns(newColumns);
    };

    const saveAndExit = () => {
        if (!cellModal) return
        setCellModal(null);

        let updatedValue: NullString;
        const rowIndex = cellModal[0];
        const col = cellModal[1];

        if (isEnumType(colType)) {
            updatedValue = { String: enumVal, Valid: enumVal !== null && enumVal !== "" };
        } else {
            switch (colType) {
                case EnumColTypes.BOOL:
                    updatedValue = { String: String(boolVal), Valid: cellVal !== null }
                    break;
                case EnumColTypes.ARRAY:
                    const hasInvalidItems = arrayItems.some(item => !item.isValid);
                    if (hasInvalidItems) {
                        return;
                    }
                    if (arrayItems.length === 0) {
                        updatedValue = { String: "", Valid: false }
                    } else {
                        const arrayValues = arrayItems.map(item =>
                            arrayType === EnumColTypes.NUMBER ? Number(item.value) : item.value
                        );
                        const arrayString = JSON.stringify(arrayValues);
                        updatedValue = { String: arrayString, Valid: true }
                    }
                    break;
                default:
                    updatedValue = {
                        String: cellVal,
                        Valid: cellVal != null && cellVal !== "" && !Number.isNaN(cellVal)
                    }
            }
        }

        if (updatedValue.Valid) {
            updateCell(updatedValue, rowIndex, col)
        }
    }

    useEffect(() => {
        const textarea = textareaRef.current;
        const numberInput = numberInputRef.current;

        if (textarea) {
            textarea.focus();
            const length = textarea.value.length;
            textarea.setSelectionRange(length, length);
        } else if (numberInput) {
            numberInput.focus();
        }

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setCellModal]);

    return (
        <div className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation} className='bg-figma-white p-6 rounded-lg'>
                {colType === EnumColTypes.TEXT && (
                    <textarea
                        ref={textareaRef}
                        className="text-figma-black mb-6 bg-figma-white border-figma-black rounded-lg w-[35rem] h-52 resize-none overflow-y-auto focus:outline-none"
                        defaultValue={cellVal}
                        onChange={(e) => setCellVal(e.target.value)}
                    />
                )}

                {colType === EnumColTypes.NUMBER && (
                    <input
                        ref={numberInputRef}
                        type="number"
                        className="text-figma-black mb-6 bg-figma-white rounded-lg w-[35rem] focus:outline-none"
                        defaultValue={cellVal}
                        onChange={(e) => setCellVal(e.target.value)}
                        placeholder="Enter a number"
                    />
                )}

                {colType === EnumColTypes.BOOL && (
                    <div className="bg-figma-white rounded-lg w-[25rem] mb-6">
                        <div className='flex flex-row justify-between items-ceter mt-4'>
                            <h2 className="text-2xl pt-0.5 mr-4">{currCol?.name}</h2>
                            <Dropdown
                                options={[{ value: "false", label: "False" },
                                { value: "true", label: "True" }]}
                                placeholder={"Select Option"}
                                defaultValue={initCellVal}
                                onSelect={(e) => setBoolVal(e.value == "true")}
                            />
                        </div>
                    </div>
                )}

                {colType === EnumColTypes.ARRAY && (
                    <div className={`bg-figma-white ${arrayItems.length == 0 && "pt-[70px]"} rounded-lg w-[35rem] max-h-96 overflow-y-auto mb-6`}>
                        <div className="">
                            {arrayItems.map((item, index) => (
                                <div key={index} className={`flex items-center gap-3 mb-3 border rounded-lg bg-figma-white p-2 ${item.isValid
                                    ? 'border-figma-gray focus:border-figma-black text-figma-black'
                                    : 'border-red-500 text-red-500'
                                    }`}>
                                    <input
                                        type="text"
                                        value={item.value}
                                        onChange={(e) => updateArrayItem(index, 'value', e.target.value)}
                                        className="grow focus:outline-none bg-figma-white"
                                        placeholder={arrayType === EnumColTypes.NUMBER ? "Enter a number" : "Enter text"}
                                    />
                                    <button
                                        onClick={() => removeArrayItem(index)}
                                        className="hover:scale-125 transition-transform duration-100"
                                    >
                                        <img src={cross} className="size-7" alt="Remove" />
                                    </button>
                                </div>
                            ))}

                            <div className="flex items-center justify-between gap-3 mt-6">
                                <button
                                    onClick={addArrayItem}
                                    className="grow p-2 border border-figma-gray hover:border-figma-black rounded-lg text-figma-black bg-figma-white transition-colors"
                                >
                                    Add Item
                                </button>
                                <Dropdown
                                    options={optionsBranches}
                                    onSelect={(option) => {
                                        const selectedType = option.value as EnumColTypes.TEXT | EnumColTypes.NUMBER;
                                        setArrayType(selectedType);
                                        setArrayItems(prev => prev.map(item => ({
                                            ...item,
                                            isValid: validateArrayItem(item.value, selectedType)
                                        })));
                                    }}
                                    placeholder={EnumColTypes.TEXT}
                                    isDown={false}
                                />
                            </div>
                        </div>
                    </div>
                )}

                {isEnumType(colType) && (
                    <div className="bg-figma-white rounded-lg w-[25rem] mb-6">
                        <div className='flex flex-row justify-between items-ceter mt-4'>
                            <h2 className="text-2xl pt-0.5 mr-4">{colType}</h2>
                            <Dropdown
                                options={enumOptions}
                                placeholder="Select value"
                                defaultValue={enumVal}
                                onSelect={(option) => setEnumVal(option.value)}
                            />
                        </div>
                    </div>
                )}

                {isConfig &&
                    <div className='flex flex-row justify-between items-ceter mt-8'>
                        <h2 className="text-2xl pt-0.5 mr-4 text-figma-black">Cell type</h2>
                        <Dropdown
                            options={[...ColTypes.map(item => ({ label: item.val, value: item.val })), 
                                     ...(currBranch?.enums || [])
                                        .filter(enumItem => enumItem.vals && enumItem.vals.length > 0)
                                        .map(enumItem => ({ label: enumItem.name, value: enumItem.name }))]}
                            defaultValue={initCell?.type.String}
                            onSelect={(val) => {
                                setCellVal("")
                                setBoolVal(false)
                                setArrayItems([])
                                setEnumVal("")
                                setColType(val.value)
                            }}
                        />
                    </div>
                }
            </div>
        </div>
    );
};

const postNewColumnData = async (col: Column, data: ColumnData, sheet: Sheet, token: string) => {
    const newColDataParams: { column: Column; data: ColumnData, sheet_id: string } = {
        column: col,
        data: data,
        sheet_id: sheet.id,
    };

    fetch(Domain + '/add_column_data', {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(newColDataParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw "Could not update column data"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

const putAdjustedColumnData = (data: ColumnData, token: string) => {
    fetch(Domain + '/update_column_data', {
        method: "put",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(data)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw "Could not update column"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default CellModal;
