import React, { useEffect, useState } from 'react';
import { useApp, ColumnProps, ColTypes } from './AppContext';
import Dropdown from './dropdown';

const ColModal = () => {
    const {
        colModal, setColModal,
        columns, setColumns,
        addColumn, setAddColumn,
    } = useApp();

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
    const [validName, setValidName] = useState(validJSON(name))

    const saveAndExit = () => {
        setColModal(-1)
        setAddColumn(false)
        if (!validName || name.length <= 0) return
        const item: ColumnProps = {
            name: name,
            columnType: columnType,
            required: required
        }
        if (addColumn) setColumns([...columns, item])
        else {
            setColumns((cols: ColumnProps[]) => {
                const newCols = [...cols];
                newCols[colModal] = item;
                return newCols;
            });
        }
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
            <div className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg min-w-[25rem] h-72 resize-none overflow-y-auto focus:outline-none"
                onClick={stopPropagation}>

                <input className={`text-figma-black text-2xl bg-figma-white mb-2 font-medium h-12 overflow-y-auto focus:outline-none
                                    ${!validName && "text-red-600"}`}
                    placeholder='Name'
                    defaultValue={name}
                    onChange={(e) => {
                        setName(e.target.value)
                        setValidName(validJSON(e.target.value) || e.target.value.length == 0)
                    }}
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
                        onSelect={(val) => {
                            if (val.value === "true") setRequired(true)
                            else setRequired(false)
                        }}
                    />
                </div>
            </div>
        </div>
    );
};

export default ColModal;
