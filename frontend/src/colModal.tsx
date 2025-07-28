import React, { useEffect, useState } from 'react';
import { useApp, ColumnProps, ColTypes, TableType, ColSuffix } from './AppContext';
import Dropdown from './dropdown';

const ColModal = () => {
    const {
        colModal, setColModal,
        columns, setColumns,
        addColumn, setAddColumn,
        setCurrTable,
        currSheet,
    } = useApp();

    const colLocalStorage = currSheet + ColSuffix
    const optionsColTypes = ColTypes.map(item => ({ label: item.val, value: item.val }));
    const [name, setName] = useState(() => {
        if (!addColumn) return columns[colModal].name
        return ''
    });
    const [columnType, setColumnType] = useState(() => {
        if (!addColumn) return columns[colModal].columnType
        return ColTypes[0].val
    });
    const [required, setRequired] = useState(() => {
        if (!addColumn) return columns[colModal].required
        return false
    });

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const validJSON = (str: string) => {
        return /^[a-zA-Z_$][a-zA-Z0-9_$\-\.]*$/.test(str)
    }

    const nameExists = (name: string) => {
        return columns.find(item => item.name == name)
    }

    const [validName, setValidName] = useState(validJSON(name))

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value)
        setValidName(
            (validJSON(e.target.value) && !nameExists(e.target.value))
            || e.target.value.length == 0
        )
    };

    const handleRequiredChange = (val: { value: string; label: string }) => {
        if (val.value === "true") setRequired(true)
        else setRequired(false)
    };

    const updateTableColumnNames = (prevName: string, newName: string) => {
        setCurrTable((prevTable: TableType) => {
            const newTable = prevTable.map(item => {
                const newItem = { ...item };
                newItem[newName] = newItem[prevName];
                delete newItem[prevName];
                return newItem;
            });
            return newTable;
        });
    };

    const updateExistingColumn = (item: ColumnProps) => {
        const prevName = columns[colModal].name
        if (prevName !== name) {
            updateTableColumnNames(prevName, name);
        }
        const newCols = [...columns];
        newCols[colModal] = item;
        setColumns(newCols);
        localStorage.setItem(colLocalStorage, JSON.stringify(newCols));
    };

    const saveAndExit = () => {
        setColModal(-1)
        setAddColumn(false)
        if (!validName || name.length <= 0) return

        const item: ColumnProps = {
            name: name,
            columnType: columnType,
            required: required
        }

        if (addColumn) {
            const newCols = [...columns, item]
            setColumns(newCols)
            localStorage.setItem(colLocalStorage, JSON.stringify(newCols));
            return
        }

        updateExistingColumn(item);
    }

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);

        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setColModal]);

    return (
        <div className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg min-w-[25rem] h-80 resize-none overflow-y-auto focus:outline-none"
                onClick={stopPropagation}>

                <input className={`text-figma-black text-2xl bg-figma-white mb-2 font-medium h-12 overflow-y-auto focus:outline-none
                                    ${(!validName && name && name != columns[colModal].name) && "text-red-600"}`}
                    placeholder='Name'
                    defaultValue={name}
                    onChange={handleNameChange}
                />

                <div className='flex flex-row justify-between items-ceter'>
                    <h2 className="text-2xl pt-0.5 mr-4">Collumn type</h2>
                    <Dropdown
                        options={optionsColTypes}
                        defaultValue={columnType}
                        onSelect={(val) => setColumnType(val.value)}
                    />
                </div>

                <div className='flex flex-row justify-between items-ceter mt-4'>
                    <h2 className="text-2xl pt-0.5 mr-4">Required</h2>
                    <Dropdown
                        options={[{ value: "false", label: "False" },
                        { value: "true", label: "True" }]}
                        defaultValue={required.toString()}
                        onSelect={handleRequiredChange}
                    />
                </div>
            </div>
        </div>
    );
};

export default ColModal;
