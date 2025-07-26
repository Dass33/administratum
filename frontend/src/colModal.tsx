import React, { useEffect, useState } from 'react';
import { useApp, ColumnProps } from './AppContext';
import Dropdown from './dropdown';

const ColModal = () => {
    const {
        setColModal,
        columns, setColumns,
    } = useApp();

    const colTypes = [
        { value: 'text', label: 'text' },
        { value: 'number', label: 'number' },
        { value: 'bool', label: 'bool' },
        { value: 'edition', label: 'edition' },
    ]

    const [name, setName] = useState('');
    const [columnType, setColumnType] = useState(colTypes[0]);
    const [required, setRequired] = useState(false);


    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const saveAndExit = () => {
        setColModal(null);
        const item: ColumnProps = {
            name: name,
            columnType: columnType.value,
            required: required
        }
        setColumns([...columns, item])
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

                <input className='text-figma-black text-2xl bg-figma-white mb-2 font-medium h-12 overflow-y-auto focus:outline-none'
                    placeholder='Name'
                    onChange={(e) => setName(e.target.value)}
                />

                <div className='flex flex-row justify-between items-ceter'>
                    <h2 className="text-2xl pt-0.5 mr-4">Collumn type</h2>
                    <Dropdown
                        options={colTypes}
                        onSelect={(val) => setColumnType(val)}
                    />
                </div>

                <div className='flex flex-row justify-between items-ceter mt-4'>
                    <h2 className="text-2xl pt-0.5 mr-4">Required</h2>
                    <Dropdown
                        options={[{ value: "false", label: "False" },
                        { value: "True", label: "True" }]}
                        onSelect={(val) => {
                            if (val.value === "True") setRequired(true)
                            else setRequired(false)
                        }}
                    />
                </div>
            </div>
        </div>
    );
};

export default ColModal;
