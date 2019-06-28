import React from "react";
import { ResponsiveContainer } from "recharts";

import ChartistGraph from "react-chartist";

import Title from "./Title";

const options = {
  fullWidth: true,
  showArea: true,
  chartPadding: {
    right: 40
  },
  high: 200,
  low: 0,
  showPoint: true,
  lineSmooth: true
};

const type = "Line";

export default ({ data }) => (
  <>
    <Title>Network Latency</Title>
    <ResponsiveContainer>
      <ChartistGraph data={data} type={type} options={options} />
    </ResponsiveContainer>
  </>
);
