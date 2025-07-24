import dropdownArrow from "./assets/dropdown_arrow.svg";

function SaveButton() {
    return (
        <div className="flex flex-row bg-figma-green text-white mr-4 rounded-lg font-bold items-center">
            <button className="hover:bg-green-700 pr-2 py-3 rounded-s-lg border-figma-white border-r-2">
                <span className="text-2xl border-figma-white ml-3">Save</span>
            </button>

            <button className="hover:bg-green-700 px-2 h-full rounded-e-lg">
                <img className="pt-1" src={dropdownArrow} alt="dropdown arrow" />
            </button>
        </div>
    );
}

export default SaveButton;
