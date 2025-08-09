import { Domain, Sheet, TableData, useApp } from "./AppContext";
import Dropdown, { DropdownOption } from "./dropdown";
import { NewItemProps } from "./NewItemModal.tsx";

type ProjectData = {
    Table: TableData,
    Sheet: Sheet,
}

function SelectProject() {
    const {
        accessToken,
        currTable, setCurrTable,
        currSheet, setCurrSheet,
        setColumns,
        setNewItemModal,
        setSheetDeleted,
        tableNames, setTableNames,
    } = useApp();
    const optionsProjects = tableNames.map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderProjectas = currTable?.name != ""
        ? currTable?.name
        : "Project"

    const setData = (data: ProjectData) => {
        setCurrTable(data.Table);
        setCurrSheet(data.Sheet);
        setColumns(data.Sheet.columns);
    }

    const assignNewName = (name: string, option: DropdownOption, setSelected: Function) => {
        renameProject(name, option.value, accessToken);

        const newTableNames = tableNames.map(idName => {
            if (idName.id === option.value) {
                setSelected({ value: idName.id, name: name })
                return { id: idName.id, name: name }
            }
            return idName;
        })
        setTableNames(newTableNames);
    }

    const deleteItem = (option: DropdownOption) => {
        deleteProject(option.value, accessToken)

        const newTableNames = tableNames.filter(
            idName => idName.id !== option.value
        )
        setTableNames(newTableNames);

        if (currTable?.id === option.value) {
            setCurrTable();
            setSheetDeleted(true);
        }
    }

    const updateValue = (option: DropdownOption, setSelected: Function) => {
        const props: NewItemProps = {
            currNames: tableNames,
            defaultIdName: { name: option.label, id: option.value },
            assignNewName(name: string) { assignNewName(name, option, setSelected) },
            deleteItem() { deleteItem(option) },
        }
        setNewItemModal(props)
    }

    return (
        <Dropdown
            options={optionsProjects}
            placeholder={placeholderProjectas}
            onSelect={(e) => getCurrTable(e.value, accessToken ?? "", setData)}
            addNewValue={(setSelected: Function) => {
                const props: NewItemProps = {
                    currNames: currSheet?.sheets_id_names ?? [],
                    assignNewName: (name) => postTable(name, accessToken ?? "", (data: ProjectData) => {
                        setData(data);
                        setSelected({ name: data.Table.name, value: data.Table.id })
                    }),
                }
                setNewItemModal(props)
            }}
            updateValue={updateValue}
        />
    );
}

const getCurrTable = (table_id: string, token: string, setData: Function) => {
    const url = Domain + `/get_project/${table_id}`;

    fetch(url, {
        method: "GET",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include"
    })
        .then(response => {
            if (response.status !== 200) {
                throw new Error("Could not retrieve project");
            }
            return response.json();
        })
        .then(result => {
            setData(result);
        })
        .catch(err => {
            console.error(err);
        });
};

const postTable = (name: string, token: string, setData: Function) => {
    const url = Domain + `/create_project`;
    const nameParam: { Name: string } = { Name: name }
    fetch(url, {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(nameParam)
    })
        .then(response => {
            if (response.status !== 201) {
                throw new Error("Could not create project");
            }
            return response.json();
        })
        .then(result => {
            setData(result);
        })
        .catch(err => {
            console.error(err);
        });
};

const renameProject = (name: string, projectId: string, token: string | undefined) => {
    const renameProjectParams: { Name: string, ProjectId: string } = {
        Name: name,
        ProjectId: projectId,
    }

    fetch(Domain + "/rename_project", {
        method: "PUT",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(renameProjectParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw new Error("Could not rename project");
            }
        })
        .catch(err => {
            console.error(err);
        });
}

const deleteProject = (projectId: string, token: string | undefined) => {
    const deleteProjectParams: { ProjectId: string } = {
        ProjectId: projectId,
    };

    fetch(Domain + '/delete_project', {
        method: "DELETE",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(deleteProjectParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                console.log(response.status)
                throw "Could not delete project"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default SelectProject;
