import React from "react";
import ChartistGraph from "react-chartist";
import { ResponsiveContainer } from "recharts";
import * as chartUtils from "../../chart-utils";
import Peak from "./Peak";
import Title from "./Title";

const MAIN_COLOR = "#D70206";

const options = {
  fullWidth: true,
  showArea: true,
  chartPadding: {
    right: 40,
  },
  high: 100,
  low: 0,
  showPoint: true,
  lineSmooth: true,
  axisX: {
    labelInterpolationFnc: chartUtils.skipLabels,
  },
};

export default ({ data }) => (
  <>
    <Title>Memory Usage (%)</Title>
    <ResponsiveContainer>
      <ChartistGraph
        data={data}
        type={"Line"}
        options={options}
        listener={chartUtils.listener(MAIN_COLOR)}
      />
    </ResponsiveContainer>
    <Peak
      value={data.series[0][data.series[0].length - 1]}
      timestamp={data.labels[data.labels.length - 1]}
      style={{ color: MAIN_COLOR }}
      unit={"%"}
    />
  </>
);
