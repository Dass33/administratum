import React, { useEffect, useRef, useState } from 'react';
import { useApp, IdName } from './AppContext';

export interface NewItemProps {
    assignNewName: (name: string) => void;
    currNames: IdName[];
    defaultIdName?: IdName;
    deleteItem?: () => void;
    ExpandingComponent?: React.FC;
}

const NewItemModal: React.FC<NewItemProps> = ({ assignNewName, currNames, defaultIdName, deleteItem, ExpandingComponent }) => {
    const {
        setNewItemModal,
    } = useApp();
    const inputRef = useRef<HTMLInputElement>(null);

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

    const [name, setName] = useState("");
    const [validName, setValidName] = useState(validJSON(name))

    const nameExists = (name: string) => {
        return currNames.find(item => item.name == name)
    }

    const saveAndExit = () => {
        setNewItemModal(null);

        if (validName && name.length > 0) {
            assignNewName(name);
        }
    }

    const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value)
        setValidName(
            (validJSON(e.target.value) && !nameExists(e.target.value))
            || e.target.value.length == 0
        )
    };

    useEffect(() => {
        const input = inputRef.current;
        if (input) {
            input.focus();
            const length = input.value.length;
            input.setSelectionRange(length, length);
        }
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, []);

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                <div className='p-4 rounded-lg bg-figma-white'>
                    <input className={`rounded-lg bg-figma-white text-figma-black text-2xl font-medium h-12
                    overflow-y-auto focus:outline-none ${(!validName && name) && "text-red-600"}`}
                        ref={inputRef}
                        placeholder='Name'
                        defaultValue={defaultIdName?.name}
                        onChange={handleNameChange}
                    />
                    {ExpandingComponent && <ExpandingComponent />}

                    {deleteItem &&
                        <div className='w-full flex justify-end mt-2 '>
                            <button className='bg-red-600 w-24 rounded-lg p-2 px-4 text-figma-white font-bold mt-4'
                                onClick={() => {
                                    deleteItem()
                                    setNewItemModal(null)
                                }}
                            >
                                <span>Delete</span>
                            </button>
                        </div>
                    }
                </div>
            </div>
        </div>
    );
};

export default NewItemModal;
