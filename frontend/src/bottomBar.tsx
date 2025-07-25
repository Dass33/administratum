import settings from "./assets/settings.svg"
import Dropdown from "./dropdown";
import { DropdownOption } from "./dropdown";

function BottomBar() {
    const optionsBranches: DropdownOption[] = [
        { value: 'config', label: 'config' },
        { value: 'questions', label: 'questions' }
    ];
    return (
        <div className="flex flex-row gap-5">
            <button className="hover:scale-110 transition-transform duration-100">
                <img className="" src={settings} />
            </button>
            <Dropdown
                options={optionsBranches}
                onSelect={() => { }}
                isDown={false}
            />
        </div>
    );
}

export default BottomBar;
