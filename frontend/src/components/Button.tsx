import { useState } from "react";

export function Button({ text, click }: { text: string; click?: () => void }) {
  const [className, setClassName] = useState("button");

  function mouseDown() {
    setClassName("button button-pressed");
  }

  function mouseUp() {
    setClassName("button");
    if (click) {
      click();
    }
  }

  function mouseLeave() {
    setClassName("button");
  }

  return (
    <button
      className={className}
      onMouseDown={mouseDown}
      onMouseUp={mouseUp}
      onMouseLeave={mouseLeave}
    >
      {text}
    </button>
  );
}
