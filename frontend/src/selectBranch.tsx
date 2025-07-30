import { useApp } from "./AppContext";
import Dropdown from "./dropdown";
import { DropdownOption } from "./dropdown";

function SelectBranch() {
    const { setBranchName } = useApp();
    const optionsBranches: DropdownOption[] = [
        { value: 'main', label: 'main' },
        { value: 'typeFix', label: 'typeFix' }
    ];
    return (
        <Dropdown
            options={optionsBranches}
            placeholder="Branch"
            onSelect={(e) => { setBranchName(e.value) }}
            addNewValue={() => { }}
        />
    );
}

export default SelectBranch;
