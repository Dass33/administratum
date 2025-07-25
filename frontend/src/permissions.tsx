import write from "./assets/write_perm.svg"
import read from "./assets/read_perm.svg"

function Permissions() {
    const icons = [write, read]

    return (
        <button className="hover:scale-125 transition-transform duration-100">
            <img src={icons[0]} />
        </button>
    );
}

export default Permissions;
