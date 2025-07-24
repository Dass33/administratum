import Dropdown from "./dropdown";
import { DropdownOption } from "./dropdown";

function SelectBranch() {
    const optionsBranches: DropdownOption[] = [
        { value: 'main', label: 'main' },
        { value: 'typeFix', label: 'typeFix' }
    ];
    return (
        <Dropdown
            options={optionsBranches}
            placeholder="Branch"
            onSelect={() => { }}
        />
    );
}

export default SelectBranch;
