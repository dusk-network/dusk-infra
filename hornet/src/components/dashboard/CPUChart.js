import React from "react";
import { makeStyles } from "@material-ui/styles";

import { ResponsiveContainer } from "recharts";

import ChartistGraph from "react-chartist";

import Title from "./Title";
import Peak from "./Peak";
import * as chartUtils from "../../chart-utils";

const MAIN_COLOR = "#0544d3";

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
  classNames: {
    line: "cpu-line",
    point: "cpu-point",
    area: "cpu-area",
  },
};

export default ({ data }) => (
  <>
    <Title>CPU Load (%)</Title>
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
