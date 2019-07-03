<template>
  <section id="audit">
    <h1>Presets</h1>
    <span :key="name" v-for="{name, sql} in presets">
      <a href="javascript:" @click="query = sql">{{ name }}</a>
      <br>
    </span>

    <h1>Query</h1>
    <textarea spellcheck="false" v-model="query"></textarea><br>
    <button @click="runQuery()">run query</button>

    <h1>Results</h1>
    <div class="container">
      <table>
        <thead>
          <tr>
            <th :key="'col'+col" v-for="col in cols">{{ col }}</th>
          </tr>
        </thead>
        <tbody>
          <tr :key="'row-'+i" v-for="(row, i) in rows">
            <td :key="'row-'+i+'-'+j" v-for="(data, j) in row" v-html="data"></td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<style>
#audit textarea {
  width: 60vw;
  height: 15em;
}
#audit .container {
  width: 100%;
  overflow: auto;
}
#audit table {
  border-collapse: collapse;
}
#audit th,
#audit td {
  padding: .2em .5em;
  border: 1px solid #999;
  white-space: nowrap;
}
</style>

<script>
"use strict";

import presets from "@/presets.js";

export default {
  data() {
    return {
      presets,
      query: "",
      cols: [],
      rows: [],
    }
  },
  methods: {
    async runQuery() {
      this.rows = [];
      this.cols = ["querying ..."];

      const resp = await fetch(`/api/database/query`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(this.query),
      });
      if(resp.status != 200){
        return alert(await resp.text());
      }

      const {cols, rows} = await resp.json();
      if(!cols.length){
        this.cols = ["no data"];
        return;
      }

      this.cols = cols;
      this.rows = rows.map(row => row.map(s => s.replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/\n/g, "<br>")));
    }
  }
}
</script>
