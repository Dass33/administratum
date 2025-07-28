import { useState } from "react";
import dropdownArrow from "./assets/dropdown_arrow.svg";
import { useApp } from "./AppContext";
import reload from "./assets/reload.svg";
import danger from "./assets/danger.svg";

type tableData = {
    currTime: number;
    projectName: string | undefined;
    branchName: string | undefined;
    length: number;
    key(index: number): string | null;
};

function SaveButton() {
    const {
        projectName,
        branchName,
    } = useApp();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const currTime = Date.now();

    const handleSubmit = async (data: tableData) => {
        setLoading(true);
        setError(null);

        try {
            const response = await fetch('http://localhost:8080/save', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    // 'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();
            console.log('Success:', result);

        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
        } finally {
            setTimeout(() => setLoading(false), 500);
        }
    };

    const currState = { ...localStorage, currTime, projectName, branchName };

    return (
        <div className="flex flex-row bg-figma-green text-white rounded-lg font-bold items-center">
            <button className="hover:bg-green-700 pr-2 py-2 rounded-s-lg border-figma-white border-r-2"
                onClick={() => handleSubmit(currState)}
            >
                {error && <img className="size-6 ml-7 mr-4 my-1" src={danger} />}
                {loading && <img className="size-6 ml-7 mr-4 my-1 animate-spin" src={reload} />}
                {!loading && !error && <span className="text-2xl border-figma-white ml-3">Save</span>}
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
