import write from "./assets/write_perm.svg"
import read from "./assets/read_perm.svg"
import { useApp, PermissionsEnum } from "./AppContext";

function Permissions() {
    const {
        currTable,
        currBranch,
    } = useApp()

    const getPermPicture = () => {
        if (!currBranch) return read
        if (!currBranch.is_protected) return write;
        if (currTable?.permision === PermissionsEnum.OWNER) {
            return write;
        }
        return read;
    }

    return (
        <button className="hover:scale-125 transition-transform duration-100">
            <img className="size-8" src={getPermPicture()} />
        </button>
    );
}

export default Permissions;
