import { useApp } from "./AppContext"
import { useEffect, useRef } from "react"

function GameView() {
    const {
        gameUrl,
        setRefresh, refresh,
    } = useApp();
    const iframeRef = useRef<HTMLIFrameElement>(null);

    useEffect(() => {
        if (!iframeRef.current || !refresh) return
        // eslint-disable-next-line no-self-assign
        iframeRef.current.src = iframeRef.current.src;
        setRefresh(false)
    }, [refresh])

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
                    </div>
                ) : (
                    <div className="embedded-container w-full h-full">
                        <div className="rounded-lg w-full h-full bg-figma-stone flex items-center justify-center">
                            <span className="text-figma-white text-2xl font-medium">
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
