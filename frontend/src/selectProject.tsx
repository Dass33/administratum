import Dropdown from "./dropdown";
import { DropdownOption } from "./dropdown";

function SelectProject() {
    const optionsProjects: DropdownOption[] = [
        { value: 'guessGame', label: 'guessGame' },
        { value: 'investingGame', label: 'investingGame' }
    ];
    return (
        <Dropdown
            options={optionsProjects}
            placeholder="Project"
            onSelect={() => { }}
        />

    );
}

export default SelectProject;
