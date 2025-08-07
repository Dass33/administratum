import { useApp } from "./AppContext"
import { useRef } from "react"

function GameView() {
    const { gameUrl } = useApp();
    const iframeRef = useRef<HTMLIFrameElement>(null);

    const handleRefresh = () => {
        if (iframeRef.current) {
            iframeRef.current.src = iframeRef.current.src;
        }
    };

    return (
        <div className="w-full h-full flex flex-col items-center justify-center py-10 pl-1">
            <div className="w-full h-full relative">
                {gameUrl.Valid ? (
                    <div className="embedded-container w-full h-full">
                        <iframe
                            ref={iframeRef}
                            className="rounded-lg w-full h-full"
                            src={gameUrl.String}
                            title="Embedded Website"
                        />
                        <button
                            onClick={handleRefresh}
                            className="absolute top-2 right-2 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                        >
                            Refresh
                        </button>
                    </div>
                ) : (
                    <div className="embedded-container w-full h-full">
                        <div className="rounded-lg w-full h-full bg-figma-stone flex items-center justify-center">
                            <span className="text-figma-white text-3xl font-medium">
                                Valid URL not provided
                            </span>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}

export default GameView;
