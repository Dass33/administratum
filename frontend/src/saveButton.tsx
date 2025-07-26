import { useState } from "react";
import dropdownArrow from "./assets/dropdown_arrow.svg";

function SaveButton() {
    return (
        <div className="flex flex-row bg-figma-green text-white rounded-lg font-bold items-center">
            <button className="hover:bg-green-700 pr-2 py-2 rounded-s-lg border-figma-white border-r-2">
                <span className="text-2xl border-figma-white ml-3">Save</span>
            </button>

            <SaveDropdown />
        </div>
    );
}

export default SaveButton;


const SaveDropdown = () => {
    const [isOpen, setIsOpen] = useState(false);
    const options = ["Export", "Share"]

    const handleSelect = (option: string): void => {
        console.log(option);
        setIsOpen(false);
    };

    return (
        <div className="relative inline-block text-left h-full">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="hover:bg-green-700 px-2 h-full rounded-e-lg"
            >
                <img className="pt-1" src={dropdownArrow} alt="dropdown arrow" />
            </button>

            {isOpen && (
                <div className="absolute z-10 my-1 w-28 right-0 bg-figma-white border border-gray-300 rounded-md shadow-lg">
                    <ul className="max-h-60 overflow-auto py-1">
                        {options.map((option, index) => (
                            <li
                                key={option || index}
                                onClick={() => handleSelect(option)}
                                className="px-4 py-2 text-sm text-gray-700 hover:bg-figma-gray hover:text-figma-black cursor-pointer"
                            >
                                {option}
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};
