import { useApp } from "./AppContext";
import refreshImg from "./assets/refresh.svg"

const RefershButton = () => {
    const { setRefresh } = useApp();

    return (
        <button className="hover:scale-125 transition-transform duration-100"
            onClick={() => setRefresh(true)}
        >
            <img className="size-8" src={refreshImg} />
        </button>
    );
}

export default RefershButton;
