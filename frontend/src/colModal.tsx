import React, { useEffect, useState } from 'react';
import { useApp, ColTypes, Column, Sheet, DEFAULT_UUID } from './AppContext';
import Dropdown from './dropdown';

type ColParams = {
    sheet_id: string
    col: Column
}

const ColModal = () => {
    const {
        colModal, setColModal,
        columns, setColumns,
        addColumn, setAddColumn,
        currSheet,
        accessToken,
    } = useApp();

    const optionsColTypes = ColTypes.map(item => ({ label: item.val, value: item.val }));
    const [name, setName] = useState(() => {
        if (!addColumn) return columns[colModal].name
        return ''
    });
    const [columnType, setColumnType] = useState(() => {
        if (!addColumn) return columns[colModal].type
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
        try {
            JSON.parse(`{"${str}": 1}`);
            return true;
        } catch (e) {
            return false;
        }
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

    const updateExistingColumn = () => {
        let newCol = columns[colModal]
        newCol.name = name;
        newCol.type = columnType
        newCol.required = required
        const newCols = [...columns];
        newCols[colModal] = newCol;
        setColumns(newCols);
        putAdjustedColumn(newCol, accessToken ?? "")
    };

    const saveAndExit = () => {
        setColModal(-1)
        if (!validName || name.length <= 0 || !currSheet) {
            return
        }
        if (addColumn) {
            const item: Column = {
                id: DEFAULT_UUID,
                name: name,
                type: columnType,
                required: required,
                data: [],
            }
            const newCols = [...columns, item]
            setColumns(newCols);
            postNewColumn(currSheet, item, accessToken ?? "");
        } else {
            updateExistingColumn();
        }
        setAddColumn(false)
    }

    const handleDeleteCol = () => {
        if (!currSheet) return;
        deleteColumn(currSheet, columns[colModal], accessToken ?? "");
        const newCols = columns.filter((_, idx) => colModal != idx);
        setColumns(newCols);
        setColModal(-1);
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
            <div className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg min-w-[25rem] resize-none overflow-y-auto focus:outline-none"
                onClick={stopPropagation}>
                <div className='h-64'>

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

                <div className='w-full flex justify-end'>
                    <button className='bg-red-600 w-24 rounded-lg p-2 px-4 text-figma-white font-bold mt-4'
                        onClick={handleDeleteCol}
                    >
                        <span>Delete</span>
                    </button>
                </div>
            </div>
        </div>
    );
};

const postNewColumn = (sheet: Sheet, col: Column, token: string) => {
    const newColParams: ColParams = {
        sheet_id: sheet.id,
        col: col,
    };

    fetch('/add_column', {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(newColParams)
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

const putAdjustedColumn = (col: Column, token: string) => {
    fetch('/update_column', {
        method: "put",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(col)
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

const deleteColumn = (sheet: Sheet, col: Column, token: string) => {
    const deleteColParams: ColParams = {
        sheet_id: sheet.id,
        col: col,
    };

    fetch('/delete_column', {
        method: "DELETE",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(deleteColParams)
    })
        .then(response => {
            if (response.status != 200) {
                throw "Could not delete column"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default ColModal;
