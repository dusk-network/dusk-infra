import React from "react";
import { ResponsiveContainer } from "recharts";

import ChartistGraph from "react-chartist";

import Title from "./Title";

const options = {};

const type = "Pie";

export default ({ data }) => (
  <>
    <Title>Disk Usage %</Title>
    <ResponsiveContainer>
      <ChartistGraph data={data} type={type} options={options} />
    </ResponsiveContainer>
  </>
);
