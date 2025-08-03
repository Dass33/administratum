import React, { useEffect, useRef, useState } from 'react';
import { useApp, EnumColTypes, Column, ColumnData, DEFAULT_UUID, Sheet, NullString } from './AppContext';
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

    const initCellVal = currCol && rowIdx < currCol.data.length
        ? currCol.data[rowIdx].value.String
        : ""

    const [cellVal, setCellVal] = useState(initCellVal);
    const [boolVal, setBoolVal] = useState(Boolean(initCellVal));
    const [arrayItems, setArrayItems] = useState<ArrayItem[]>([]);
    const [arrayType, setArrayType] = useState<EnumColTypes.TEXT | EnumColTypes.NUMBER>(EnumColTypes.TEXT);

    const colType = cellModal
        ? cellModal[1].type
        : EnumColTypes.TEXT;

    const optionsBranches: DropdownOption[] = [
        { value: EnumColTypes.TEXT, label: EnumColTypes.TEXT },
        { value: EnumColTypes.NUMBER, label: EnumColTypes.NUMBER }
    ];

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
        let newCol = { ...col };

        newCol.data = col.data.map(item => {
            if (item.idx === rowIndex) {
                itemFound = true;
                item_id = item.id;
                return {
                    id: item_id,
                    idx: rowIndex,
                    value: newVal,
                };
            }
            return item;
        });

        if (itemFound) {
            const updatedData: ColumnData = {
                id: item_id,
                idx: rowIndex,
                value: newVal,
            };
            putAdjustedColumnData(updatedData, accessToken ?? "");
        } else {
            const newColData: ColumnData = {
                id: item_id,
                idx: rowIndex,
                value: newVal,
            };
            newCol.data.push(newColData);
            postNewColumnData(col, newColData, currSheet, accessToken ?? "");

            let newSheet = currSheet;
            if (newSheet && newSheet.row_count - 1 <= rowIndex) newSheet.row_count++;
            setCurrSheet(newSheet);
        }

        const newColumns = columns.map(item => {
            if (item.id === col.id) return newCol;
            return item;
        });
        setColumns(newColumns);
    };

    const removeEmptyRow = (newVal: any, rowIndex: number): boolean => {
        // if (newVal) return false;
        //
        // if (currTable.length - 1 == rowIndex) return false
        // const len = Object.entries(currTable[rowIndex]).filter(([key, val]) => {
        //     return key != col && val
        // }).length
        // if (!len && (newVal === null || newVal === "" || newVal === undefined)) {
        //     setCurrTable(currTable.filter((_, idx) => { return rowIndex != idx }))
        //     return true
        // }
        return false
    }

    const saveAndExit = () => {
        if (!cellModal) return
        setCellModal(null);

        let updatedValue: NullString;
        const rowIndex = cellModal[0];
        const col = cellModal[1];

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
                    updatedValue = { String: arrayString, Valid: false }
                }
                break;
            default:
                updatedValue = { String: cellVal, Valid: cellVal !== null }
        }

        if (!removeEmptyRow(updatedValue, rowIndex) && updatedValue) {
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
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                {colType === EnumColTypes.TEXT && (
                    <textarea
                        ref={textareaRef}
                        className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg w-[35rem] h-72 resize-none overflow-y-auto focus:outline-none"
                        defaultValue={cellVal}
                        onChange={(e) => setCellVal(e.target.value)}
                    />
                )}

                {colType === EnumColTypes.NUMBER && (
                    <input
                        ref={numberInputRef}
                        type="number"
                        className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg w-[35rem] focus:outline-none"
                        defaultValue={cellVal}
                        onChange={(e) => setCellVal(e.target.value)}
                        placeholder="Enter a number"
                    />
                )}

                {colType === EnumColTypes.BOOL && (
                    <div className="bg-figma-white p-6 rounded-lg w-[35rem] mb-6">
                        <label className="flex items-center cursor-pointer">
                            <input
                                type="checkbox"
                                checked={boolVal}
                                onChange={(e) => setBoolVal(e.target.checked)}
                                className="sr-only"
                            />
                            <div className={`relative w-12 h-6 rounded-full transition-colors duration-200 ${boolVal ? 'bg-figma-black' : 'bg-figma-gray'
                                }`}>
                                <div className={`absolute top-1 left-1 w-4 h-4 bg-figma-white rounded-full transition-transform duration-200 ${boolVal ? 'translate-x-6' : 'translate-x-0'
                                    }`} />
                            </div>
                            <span className="ml-3 text-figma-black">
                                {boolVal ? 'True' : 'False'}
                            </span>
                        </label>
                    </div>
                )}

                {colType === EnumColTypes.ARRAY && (
                    <div className={`bg-figma-white p-7 ${arrayItems.length == 0 && "pt-[70px]"} rounded-lg w-[35rem] max-h-96 overflow-y-auto mb-6`}>
                        <div className="">
                            {arrayItems.map((item, index) => (
                                <div key={index} className="flex items-center gap-3 mb-3">
                                    <input
                                        type="text"
                                        value={item.value}
                                        onChange={(e) => updateArrayItem(index, 'value', e.target.value)}
                                        className={`flex-1 p-2 rounded-lg border focus:outline-none ${item.isValid
                                            ? 'bg-figma-white border-figma-gray focus:border-figma-black text-figma-black'
                                            : 'bg-figma-white border-red-500 text-red-500'
                                            }`}
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

                            <div className="flex items-center gap-3 mt-6">
                                <button
                                    onClick={addArrayItem}
                                    className="flex-1 p-2 border border-figma-gray hover:border-figma-black rounded-lg text-figma-black bg-figma-white transition-colors"
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
                                <div className="size-7"></div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

const postNewColumnData = async (col: Column, data: ColumnData, sheet: Sheet, token: string) => {
    const newColDataParams: { column: Column; data: ColumnData, sheet_id: String } = {
        column: col,
        data: data,
        sheet_id: sheet.id,
    };

    fetch('/add_column_data', {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(newColDataParams)
    })
        .then(response => {
            if (response.status != 200) {
                throw "Could not update column data"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

const putAdjustedColumnData = (data: ColumnData, token: string) => {
    fetch('/update_column_data', {
        method: "put",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(data)
    })
        .then(response => {
            if (response.status != 200) {
                throw "Could not update column"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default CellModal;
