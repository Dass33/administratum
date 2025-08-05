import { useState } from "react";
import dropdownArrow from "./assets/dropdown_arrow.svg";
import { useApp, Domain } from "./AppContext";

function SaveButton() {
    const {
        setShareModal,
    } = useApp();

    return (
        <div className="flex flex-row bg-figma-green text-white rounded-lg font-bold items-center">
            <button className="hover:bg-green-700 pr-2 py-2 rounded-s-lg border-figma-white border-r-2"
                onClick={() => setShareModal(true)}
            >
                <span className="text-2xl border-figma-white ml-3">Share</span>
            </button>

            <SaveDropdown />
        </div>
    );
}

export default SaveButton;


enum SaveOptoins {
    EXPORT = "Export",
    COPY_LINK = "Copy Link",
}

const SaveDropdown = () => {
    const {
        currSheet,
    } = useApp();
    const [isOpen, setIsOpen] = useState(false);
    const options = [SaveOptoins.EXPORT, SaveOptoins.COPY_LINK]

    const handleSelect = (option: string): void => {
        switch (option) {
            case SaveOptoins.EXPORT:
                break;
            case SaveOptoins.COPY_LINK:
                navigator.clipboard.writeText(`${Domain}/json/${currSheet?.branch_id_name.id}`);
        }
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
                <div className="absolute z-10 my-1 w-[123px] right-0 bg-figma-white border border-gray-300 rounded-md shadow-lg">
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
