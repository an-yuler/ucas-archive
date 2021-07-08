<template>
  <div class="q1">
    <div class="query-form">
      <el-alert
        title="任务二"
        :closable="false"
        type="info"
        description="实现多层股权的穿透式查询，可以根据指定层数获得对应层级的股东，例如：输入“招商局轮船股份有限公司”和层数3"
      >
      </el-alert>
      <el-form :rules="rules" ref="form" :model="form" label-width="80px">
        <el-form-item label="公司名称" prop="company">
          <el-col :span="7">
            <el-input
              v-model="form.company"
              placeholder="请输入公司名"
            ></el-input>
          </el-col>
        </el-form-item>
        <el-form-item label="跳数" prop="hop">
          <el-col :span="2">
            <el-input
              v-model.number="form.hop"
              placeholder="数字1-20"
            ></el-input>
          </el-col>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit('form')">查询</el-button>
          <el-button type="info" @click="onUseExample">示例</el-button>
          <el-button @click="onClearForm('form')">清空</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="res-show">
      <el-table
        :data="responseVO.results"
        border
        height="700"
        style="width: 100%"
      >
        <el-table-column type="index" width="80" label="id"> </el-table-column>
        <el-table-column prop="chain" label="chain">
          <template slot-scope="scope">
            {{ scope.row.chain.join("→") }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import Api from "@/api/index.js";
export default {
  data() {
    var checkHop = (rule, value, callback) => {
      if (!value) {
        return callback(new Error("必填项"));
      }
      setTimeout(() => {
        if (!Number.isInteger(value)) {
          callback(new Error("值应该为数字"));
        } else {
          if (value > 20 || value < 1) {
            callback(new Error("建议值为(1-20)"));
          } else {
            callback();
          }
        }
      }, 1000);
    };
    return {
      form: {
        company: "",
        hop: "",
      },
      responseVO: "",
      rules: {
        company: [{ required: true, message: "必填项", trigger: "blur" }],
        hop: [{ required: true, validator: checkHop, trigger: "blur" }],
      },
    };
  },
  methods: {
    onSubmit(formName) {
      this.$refs[formName].validate((valid) => {
        var vo = { params: [this.form.company, this.form.hop.toString()] };
        if (valid) {
          Api.query2(vo).then((response) => {
            if (response.data.code != 100) {
              this.$message({
                showClose: true,
                message: "查询出错",
                type: "error",
              });
              return;
            }

            this.responseVO = response.data;
            this.$message({
              showClose: true,
              message: "查询完成",
              type: "success",
            });
          });
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    onClearForm(formName) {
      this.$refs[formName].resetFields();
      this.responseVO = "";
    },
    onUseExample() {
      this.form.company = "招商局轮船股份有限公司";
      this.form.hop = 3;
    },
  },
};
</script>
<style scoped>
.el-alert {
  margin: 10px 0;
}

.res-show {
  margin-top: 100px;
}
</style>
