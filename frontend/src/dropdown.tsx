import React, { useState, useEffect } from 'react';
import plus from "./assets/plus.svg";
import settings from "./assets/settings.svg";

export interface DropdownOption {
    value: string;
    label: string;
}

export interface DropdownProps {
    options: DropdownOption[];
    placeholder?: string;
    onSelect?: (option: DropdownOption) => void;
    isDown?: boolean;
    defaultValue?: string;
    addNewValue?: () => void;
    updateValue?: (option: DropdownOption) => void;
}

const Dropdown: React.FC<DropdownProps> = ({
    options,
    placeholder = "Select an option",
    onSelect,
    isDown = true,
    defaultValue,
    addNewValue,
    updateValue,
}) => {
    const [isOpen, setIsOpen] = useState(false);
    const [selectedOption, setSelectedOption] = useState<DropdownOption | undefined>(() => {
        if (defaultValue) return options.find((item) => item.value === defaultValue)
        if (options.length > 0) return options[0]
        return undefined
    });

    useEffect(() => {
        const selected = options.find((item) => item.value === selectedOption?.value);
        if (selected && selected.label !== selectedOption?.label) {
            setSelectedOption(selected);
        }
    }, [options]);

    const handleSelect = (option: DropdownOption): void => {
        setSelectedOption(option);
        setIsOpen(false);
        if (onSelect) {
            onSelect(option);
        }
    };

    return (
        <div className="relative inline-block text-left min-w-28">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="w-full bg-figma-white border border-figma-gray rounded-lg px-4 py-2 text-left shadow-sm hover:bg-gray-50 focus:border-figma-black pr-8"
            >
                <span className="block truncate">
                    {selectedOption?.label || placeholder}
                </span>
                <span className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                    <svg
                        className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''
                            }`}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                </span>
            </button>

            {isOpen && (
                <div className={`absolute z-10 mt-1 w-full bg-figma-white border border-gray-300 rounded-md shadow-lg ${isDown ? 'top-full' : 'bottom-full mb-1'
                    }`}>
                    <ul className="max-h-60 overflow-auto py-1">
                        {options.map((option, index) => (
                            <div className='flex items-center hover:bg-figma-gray hover:text-figma-black justify-between pl-4 pr-2 py-2'
                                key={option.value || index}
                            >
                                <li
                                    onClick={() => handleSelect(option)}
                                    className='mr-2 text-sm text-gray-700 cursor-pointer h-5 overflow-hidden text-ellipsis whitespace-nowrap flex-1'
                                >
                                    {option.label}
                                </li>
                                {updateValue &&
                                    <button className='w-4 h-4 hover:scale-110 transition-transform duration-100 flex-shrink-0'
                                        onClick={() => {
                                            setIsOpen(false);
                                            updateValue(option);
                                        }}
                                    >
                                        <img className="w-full h-full" src={settings} />
                                    </button>
                                }
                            </div>
                        ))}
                        {addNewValue &&
                            <li
                                onClick={() => {
                                    setIsOpen(false);
                                    addNewValue();
                                }}
                                className="px-4 py-2 text-sm text-gray-700 hover:bg-figma-gray hover:text-figma-black cursor-pointer flex"
                            >
                                <span className='font-medium'>New</span>
                                <img className='size-5 ml-1' src={plus} />
                            </li>
                        }
                    </ul>
                </div>
            )}
        </div>
    );
};

export default Dropdown;
