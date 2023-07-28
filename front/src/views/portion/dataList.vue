<script setup>
import {ref} from "vue";

const small = ref(false)
const background = ref(false)
const disabled = ref(false)

</script>

<template>
  <el-scrollbar >
    <div class="back-icon" @click="closeDetail">
      <el-icon >
        <Back />
      </el-icon>
    </div>


    <el-table class="view-data" :data="tableDataList">

      <el-table-column
          show-overflow-tooltip

          v-for="column in columns"
          :label="column"
          :key="column"
          :prop="column"
      >

      </el-table-column>

    </el-table>


    <div class="demo-pagination-block">
      <el-pagination
          v-model:current-page="currentPage"
          :page-size="currentSize"
          :small="small"
          :disabled="disabled"
          :background="background"
          layout="total, prev, pager, next"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
      />
    </div>

  </el-scrollbar>
</template>

<style scoped>
.back-icon {
  font-size: 25px
}
</style>

<script>

import requestFunc from "../../request/table";
import {useCounterStore} from "../../stores/stores";

export default {
  data() {
    return {
      tableStruct: [],
      tableDataList: [],
      columns: [],

      currentPage: 1,
      currentSize: 15,
      total: 15,
    }
  },
  methods: {
    getDetail(dp) {

      // console.log("useCounterStore().getDataPage()", JSON.stringify(dp) )
      requestFunc.GetDataList(dp).then(result => {
        // console.log("test",JSON.stringify(result))
        this.tableDataList = result.list
        this.total = result.total
        for (const index in result.list) {
          this.columns = Object.keys(result.list[index])
        }
        // console.log("test",JSON.stringify(this.columns))
        this.$message.success('请求成功');
      }).catch(error => {
        console.log("error: ",error)
        this.$message.error('请求失败');
      });
    },

    getDataPage(table_name, page) {
      const dp = useCounterStore().getDataPage()
      if (table_name !== "") {
        dp.table_name = table_name
      }
      dp.page = page
      dp.size = this.currentSize
      useCounterStore().setDataPage(dp)
      // console.log(JSON.stringify(dp))
      return dp
    },

    handleCurrentChange(val) {
      const dp = this.getDataPage("", val)
      this.getDetail(dp)
    },

    handleSizeChange(){
      console.log("page-size 改变时触发")
    },

    closeDetail(){
      this.$emit('closeDataList-flag');
    },

    dataList(){
      const table_name = useCounterStore().getTableName()
      const dp = this.getDataPage(table_name, this.currentPage)
      this.getDetail(dp)
    },
  },

  mounted() {
    // 在其他方法或是生命周期中也可以调用方法
    this.dataList()
  }
}

</script>
