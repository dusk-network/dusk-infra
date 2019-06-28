import React from "react";

import Title from "./Title";
import {
  BarChart,
  Bar,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from "recharts";

export default ({ data }) => (
  <>
    <Title>Highest Score</Title>
    <ResponsiveContainer>
      <BarChart width={150} height={40} data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="height" />
        <YAxis dataKey="score" />
        <Bar dataKey="score" fill="#dd416a" />
      </BarChart>
    </ResponsiveContainer>
  </>
);
