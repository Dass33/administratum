import { useApp } from "./AppContext";
import Dropdown from "./dropdown";
import { DropdownOption } from "./dropdown";

function SelectProject() {
    const { setProjectName } = useApp();
    const optionsProjects: DropdownOption[] = [
        { value: 'guessGame', label: 'guessGame' },
        { value: 'investingGame', label: 'investingGame' }
    ];
    return (
        <Dropdown
            options={optionsProjects}
            placeholder="Project"
            onSelect={(e) => { setProjectName(e.value) }}
            addNewValue={() => { }}
        />

    );
}

export default SelectProject;
