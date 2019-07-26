import { createSelector } from "reselect";

const lastBlockInfo = state => state.lastBlock;
const timeSelector = state => state.blockTime || [];
const txSelector = state => state.txs || [];
const cpuSelector = state => state.cpu || [];
const logSelector = state => state.log || [];
const netSelector = state => state.net || [];
const memSelector = state => state.memory || [];
const diskSelector = state => state.disk || [];
const warnSelector = state => state.warnings || [];
const threadSelector = state => state.thread || [];

export const getCurrentBlockInfo = createSelector(
  lastBlockInfo,
  lastBlock => lastBlock
);

export const getWarnings = createSelector(
  warnSelector,
  warnings => warnings.slice(0, 200)
);

export const getTimeMetrics = createSelector(
  timeSelector,
  timest =>
    timest
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value.toFixed(2));
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getTxMetrics = createSelector(
  txSelector,
  txs =>
    txs
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getCPUMetrics = createSelector(
  cpuSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value.toFixed(1));
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getNetMetrics = createSelector(
  netSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getLogMetrics = createSelector(
  logSelector,
  info => info.slice(0, 200)
);

export const getMemoryMetrics = createSelector(
  memSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getThreadMetrics = createSelector(
  threadSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getDiskMetrics = createSelector(
  diskSelector,
  info => {
    if (info[0]) {
      let used = info[0].value;
      let free = 100 - used;
      return {
        series: [used, free],
        labels: [`Used: ${used}%`, `Free: ${free}%`]
      };
    }
  }
);
