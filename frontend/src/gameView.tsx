import { useApp } from "./AppContext"
import expand from "./assets/expand.svg"

function GameView() {
    const { gameUrl } = useApp();

    return (
        <div className="flex items-center">
            <div className="flex flex-row">
                <button className="size-8 hover:scale-125 transition-transform duration-100" >
                    <img
                        src={expand}
                        alt="dropdown arrow"
                    />
                </button>
                <div className="embedded-container">
                    <iframe className="rounded-lg"
                        src={gameUrl}
                        width="370"
                        height="750"
                        title="Embedded Website"
                    />
                </div>
            </div>
        </div>
    )
}

export default GameView
