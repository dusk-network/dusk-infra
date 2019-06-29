/* eslint-disable no-script-url */

import { withStyles } from "@material-ui/core";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import { makeStyles } from "@material-ui/styles";
import React from "react";

const cellMap = {
  warning:withStyles(theme => ({
    color: "#f05b4f"
  }))(TableCell),

  error:withStyles(theme => ({
    color: "#d70206"
  }))(TableCell),

  fatal:withStyles(theme => ({
    color: "#6b0392"
  }))(TableCell),

  panic : withStyles(theme => ({
    color: "#0544d3"
  }))(TableCell)
}


const useStyles = makeStyles(theme => ({}));

export default function Nodes({ items = [] }) {
  const classes = useStyles();
  const Cell = cellMap["error"]

  return (
    <>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Timestamp</TableCell>
            <TableCell>Error</TableCell>
            <TableCell>Message</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map(({ timestamp, warnings: { level, error="-", msg } }, index) => {
            // let TR = trMap[level];
            return <>
            <TableRow key={`${timestamp}-${index}`} >
              <Cell>{timestamp}</Cell>
              <Cell>{error}</Cell>
              <Cell>{msg}</Cell>
            </TableRow>
            </>
          })}
        </TableBody>
      </Table>
    </>
  );
}
