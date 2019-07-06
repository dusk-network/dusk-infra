import Chartist from "chartist";

const lz = (num, length = 2) => String(num).padStart(length, "0");
const step = (index, length) => {
  if (index > 0 && index < length - 1 && index % Math.round(length / 3) === 0) {
    return true;
  }
  return false;
};

export const listener = className => ({
  draw(data) {
    if (data.type === "point" && step(data.index, data.series.length)) {
      data.group.append(
        new Chartist.Svg(
          "circle",
          {
            cx: data.x,
            cy: data.y,
            r: 4,
          },
          className
        )
      );
    }
  },
});

export function skipLabels(value, index, labels) {
  if (step(index, labels.length)) {
    const date = new Date(value);
    return `${lz(date.getHours())}:${lz(date.getMinutes())}:${lz(
      date.getSeconds()
    )}`;
  }
}
