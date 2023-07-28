<script setup>

import DataList from "./dataList.vue";
import Form from "./form.vue";

</script>

<template>
  <!--   内容   -->
  <el-scrollbar v-if="showStruct" style="height: 100%">
    <el-table :data="tableStruct">
      <el-table-column prop="number" label="序号" width="140"/>
      <el-table-column prop="name" label="表名" width="120"/>

      <el-table-column label="字段">
        <template #default="{ row }">
          <div v-for="(value, key) in row.fields" :key="key">
            {{ key }}
          </div>
        </template>
      </el-table-column>

      <el-table-column label="类型">
        <template #default="{ row }">
          <div v-for="(value, key) in row.fields" :key="key">
            {{ value }}
          </div>
        </template>
      </el-table-column>


      <el-table-column fixed="right" label="操作" width="120">
        <template #default="{ row }">
          <el-button link type="primary" size="large" @click="openDetail(row.name)">Detail</el-button>
          <el-button link type="primary" size="large" @click="openForm(row.fields)">Edit</el-button>
        </template>
      </el-table-column>

    </el-table>

  </el-scrollbar>

  <DataList v-if="showDataList" @closeDataList-flag="closeDataList"></DataList>

  <Form v-if="showDataForm" @closeForm-flag="closeForm"></Form>
</template>

<style scoped>


</style>


<script>

import requestFunc from "../../request/table";
import {useCounterStore} from '../../stores/stores';

export default {
  data() {
    return {
      showStruct: true,
      showDataList: false,
      showDataForm: false,
      tableStruct: [],

    }
  },

  methods: {
    getTables() {
      // console.log("useCounterStore().getDbType()",useCounterStore().getDbType())
      requestFunc.GetTables(useCounterStore().getDbType()).then(result => {
        this.tableStruct = result.tableStruct
        // console.log(JSON.stringify(this.tableStruct))
        this.$message.success('请求成功');
      }).catch(error => {
        this.$message.error('请求失败');
      });
    },



    openDetail(table_name) {
      this.showStruct = false // 关闭表结构列表，显示详细数据内容
      this.showDataList = true
      useCounterStore().setTableName(table_name)
    },

    closeDataList(){
      this.showStruct = !this.showStruct // 关闭详细数据内容，显示表结构列表
      this.showDataList = !this.showDataList
    },


    openForm(data_struct){
      this.showStruct = false
      this.showDataForm = true
      useCounterStore().setDataStruct(data_struct)

    },

    closeForm(){
      this.showStruct = !this.showStruct
      this.showDataForm = !this.showDataForm
    },

  },
  mounted() {
    // 在其他方法或是生命周期中也可以调用方法
    this.getTables()
  }
}


</script>
