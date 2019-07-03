import React from "react";
import { ResponsiveContainer } from "recharts";

import ChartistGraph from "react-chartist";

import Title from "./Title";

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
    showLabel: false,
  },
};

const type = "Line";

export default ({ data }) => (
  <>
    <Title>Memory Usage</Title>
    <ResponsiveContainer>
      <ChartistGraph data={data} type={type} options={options} />
    </ResponsiveContainer>
  </>
);
