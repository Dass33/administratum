import React, { useState, useEffect, useRef } from 'react';
import plus from "./assets/plus.svg";
import menu from "./assets/menu.svg";

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
    addNewValue?: (setSelected: (option: DropdownOption) => void) => void;
    updateValue?: (option: DropdownOption, setSelected: (option: DropdownOption) => void) => void;
    EveryRender?: (setSelected: (option: DropdownOption) => void, currentSelected?: DropdownOption) => void;
    usePortal?: boolean;
}

const Dropdown: React.FC<DropdownProps> = ({
    options,
    placeholder = "Select an option",
    onSelect,
    isDown = true,
    defaultValue,
    addNewValue,
    updateValue,
    EveryRender,
    usePortal = false,
}) => {
    const [isOpen, setIsOpen] = useState(false);
    const [selectedOption, setSelectedOption] = useState<DropdownOption | undefined>(() => {
        if (defaultValue) return options.find((item) => item.value === defaultValue)
        if (options.length > 0) return options[0]
        return undefined
    });
    const [dropdownPosition, setDropdownPosition] = useState({ top: 0, left: 0, width: 0 });
    const buttonRef = useRef<HTMLButtonElement>(null);

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

    const updateDropdownPosition = () => {
        if (buttonRef.current && usePortal) {
            const rect = buttonRef.current.getBoundingClientRect();
            setDropdownPosition({
                top: isDown ? rect.bottom + window.scrollY : rect.top + window.scrollY,
                left: rect.left + window.scrollX,
                width: rect.width
            });
        }
    };

    useEffect(() => {
        if (isOpen && usePortal) {
            updateDropdownPosition();
            const handleResize = () => updateDropdownPosition();
            window.addEventListener('resize', handleResize);
            return () => window.removeEventListener('resize', handleResize);
        }
    }, [isOpen, usePortal, isDown]);

    if (EveryRender) EveryRender(setSelectedOption, selectedOption);

    return (
        <div className="relative inline-block text-left min-w-28">
            <button
                ref={buttonRef}
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
                usePortal ? (
                    <div
                        className="fixed z-[60] bg-figma-white border border-gray-300 rounded-md shadow-lg"
                        style={{
                            top: isDown ? dropdownPosition.top + 4 : dropdownPosition.top - 4,
                            left: dropdownPosition.left,
                            width: dropdownPosition.width,
                            transform: isDown ? 'none' : 'translateY(-100%)'
                        }}
                    >
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
                                        <button className='w-5 pl-0.5 hover:scale-110 transition-transform duration-100 flex-shrink-0'
                                            onClick={() => {
                                                setIsOpen(false);
                                                updateValue(option, setSelectedOption);
                                            }}
                                        >
                                            <img className="w-full h-full" src={menu} />
                                        </button>
                                    }
                                </div>
                            ))}
                            {addNewValue &&
                                <li
                                    onClick={() => {
                                        setIsOpen(false);
                                        addNewValue(setSelectedOption);
                                    }}
                                    className="px-4 py-2 text-sm text-gray-700 hover:bg-figma-gray hover:text-figma-black cursor-pointer flex"
                                >
                                    <span className='font-medium'>New</span>
                                    <img className='size-5 ml-1' src={plus} />
                                </li>
                            }
                        </ul>
                    </div>
                ) : (
                    <div className={`absolute z-[60] mt-1 w-full bg-figma-white border border-gray-300 rounded-md shadow-lg ${isDown ? 'top-full' : 'bottom-full mb-1'
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
                                        <button className='w-5 pl-0.5 hover:scale-110 transition-transform duration-100 flex-shrink-0'
                                            onClick={() => {
                                                setIsOpen(false);
                                                updateValue(option, setSelectedOption);
                                            }}
                                        >
                                            <img className="w-full h-full" src={menu} />
                                        </button>
                                    }
                                </div>
                            ))}
                            {addNewValue &&
                                <li
                                    onClick={() => {
                                        setIsOpen(false);
                                        addNewValue(setSelectedOption);
                                    }}
                                    className="px-4 py-2 text-sm text-gray-700 hover:bg-figma-gray hover:text-figma-black cursor-pointer flex"
                                >
                                    <span className='font-medium'>New</span>
                                    <img className='size-5 ml-1' src={plus} />
                                </li>
                            }
                        </ul>
                    </div>
                )
            )}
        </div>
    );
};

export default Dropdown;
