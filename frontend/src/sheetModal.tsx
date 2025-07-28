import React, { useEffect, useRef, useState } from 'react';
import { Sheets, useApp } from './AppContext';

const SheetModal = () => {
    const {
        setSheetModal,
        sheets, setSheets
    } = useApp();
    const textareaRef = useRef<HTMLTextAreaElement>(null);

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const [name, setName] = useState("");

    const validJSON = (str: string) => {
        return /^[a-zA-Z_$][a-zA-Z0-9_$\-\.]*$/.test(str)
    }

    const nameExists = (name: string) => {
        return sheets.find(item => item == name)
    }

    const [validName, setValidName] = useState(validJSON(name))

    const saveAndExit = () => {
        setSheetModal(false);

        if (validName && name.length > 0) {
            const newSheets = [...sheets, name];
            setSheets(newSheets);
            localStorage.setItem(Sheets, JSON.stringify(newSheets));
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
        const textarea = textareaRef.current;

        if (textarea) {
            textarea.focus();
            const length = textarea.value.length;
            textarea.setSelectionRange(length, length);
        }
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setSheetModal]);

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                <input className={`p-4 rounded-lg text-figma-black text-2xl bg-figma-white mb-2 font-medium h-12
                    overflow-y-auto focus:outline-none ${(!validName && name) && "text-red-600"}`}
                    placeholder='Name'
                    onChange={handleNameChange}
                />
            </div>
        </div>
    );
};

export default SheetModal;
