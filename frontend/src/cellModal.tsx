import React, { useEffect, useRef, useState } from 'react';
import { useApp, TableType, EnumColTypes } from './AppContext';
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
        currTable, setCurrTable,
    } = useApp();
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    const numberInputRef = useRef<HTMLInputElement>(null);

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const initCellVal = cellModal
        ? currTable[cellModal[0]][cellModal[1].name]
        : ""

    const [cellVal, setCellVal] = useState(initCellVal);
    const [boolVal, setBoolVal] = useState(Boolean(initCellVal));
    const [arrayItems, setArrayItems] = useState<ArrayItem[]>([]);
    const [arrayType, setArrayType] = useState<EnumColTypes.TEXT | EnumColTypes.NUMBER>(EnumColTypes.TEXT);

    const colType = cellModal
        ? cellModal[1].columnType
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

    const updateCell = (newVal: any, rowIndex: number, col: string) => {
        setCurrTable((prevTable: TableType) => {
            const newTable = [...prevTable];
            newTable[rowIndex][col] = newVal;
            return newTable;
        });
    }

    const removeEmptyRow = (newVal: any, rowIndex: number, col: string): boolean => {
        if (currTable.length - 1 == rowIndex) return false
        const len = Object.entries(currTable[rowIndex]).filter(([key, val]) => {
            return key != col && val
        }).length
        if (!len && !newVal) {
            setCurrTable(currTable.filter((_, idx) => { return rowIndex != idx }))
            return true
        }
        return false
    }

    const saveAndExit = () => {
        setCellModal(null);
        if (!cellModal) return

        let updatedValue: any;
        const rowIndex = cellModal[0];
        const col = cellModal[1];

        switch (colType) {
            case EnumColTypes.TEXT:
                updatedValue = cellVal;
                break;
            case EnumColTypes.NUMBER:
                updatedValue = cellVal === '' ? null : Number(cellVal);
                break;
            case EnumColTypes.BOOL:
                updatedValue = boolVal;
                break;
            case EnumColTypes.ARRAY:
                const hasInvalidItems = arrayItems.some(item => !item.isValid);
                if (hasInvalidItems) {
                    return;
                }

                if (arrayItems.length === 0) {
                    updatedValue = null;
                } else {
                    const arrayValues = arrayItems.map(item =>
                        arrayType === EnumColTypes.NUMBER ? Number(item.value) : item.value
                    );
                    updatedValue = JSON.stringify(arrayValues);
                }
                break;
            default:
                updatedValue = cellVal;
        }

        if (!removeEmptyRow(updatedValue, rowIndex, col.name)) {
            updateCell(updatedValue, rowIndex, col.name)
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
                        className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg w-[35rem] border border-figma-gray focus:border-figma-black focus:outline-none"
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

export default CellModal;
