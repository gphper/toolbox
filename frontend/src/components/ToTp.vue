<script setup>
import {reactive,onMounted} from 'vue'
import {ScreenShot,Totp,Storage,Get} from '../../wailsjs/go/main/App'
import jsQR from 'jsqr';
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const centerDialogVisible = ref(false)
const addFormVisible = ref(false)
const cropper = ref(null)
const formRef = ref(null)
let reloadTimeHandle = null

const data = reactive({
  img:"",
  qrcode:"",
  username:"匿名用户",
  code:"------",
  qdata:"",
})

const form = reactive({
  account:'',
  secret:'',
})

const rules = reactive({
  account:[{ required: true, message: '请输入账号', trigger: 'blur' }],
  secret: [{ required: true, message: '请输入秘钥', trigger: 'blur' }],
})

// 提交表单
function submitForm(){
  formRef.value.validate((valid) => {
    console.log(valid)
  })
  
  // 组装数据
  let data2Fa = "otpauth://totp/GitHub:"+form.account+"?secret="+form.secret
  data.qdata = data2Fa
  Storage(data2Fa).then().catch(err => {
    //TODO
  })
  // 生成验证码
  reloadTotp()
  addFormVisible.value = false
}

// 显示表单
function showForm(){
  addFormVisible.value = true
  form.account = ''
  form.secret = ''
}

// 截屏
function screenShot(){
  window.runtime.WindowHide()
  setTimeout(function(){
    ScreenShot().then(result => {
      data.img = result
      window.runtime.WindowShow()
      centerDialogVisible.value = true
    }).catch(err => {
      ElMessage({
        message: '截图失败：' + err,
        type: 'warning',
      })
    })
  },300)
}

// 截屏操作
function screenCut(){
  cropper.value.getCropData(result => {
    data.qrcode = result
    decodeQr()
  })
}

// 设置剪切板
function copy(){
  let s = window.runtime.ClipboardSetText(data.code)
  if (s) {
      ElMessage({
        message: '复制成功',
        type: 'success',
      })
  }else{
    ElMessage({
        message: '复制失败,请手动复制',
        type: 'warning',
      })
  }
}

// 解码二维码
function decodeQr(){
  
  const image = new Image();
  image.src = `${data.qrcode}`;

  image.onload = () => {
      const canvas = document.createElement('canvas');
      const context = canvas.getContext('2d');
      canvas.width = image.width;
      canvas.height = image.height;
      context.drawImage(image, 0, 0, image.width, image.height);
      const imageData = context.getImageData(0, 0, image.width, image.height);
      const code = jsQR(imageData.data, imageData.width, imageData.height);
      if (code) {
        data.qdata = code.data
        // 在此处理识别到的二维码结果
        ElMessage({
          message: '识别到二维码内容：' + code.data,
          type: 'success',
        })
        // 开始更新验证码
        reloadTotp();
        // 隐藏弹窗
        centerDialogVisible.value = false
        Storage(data.qdata)
      } else {
        ElMessage({
          message: '未识别到二维码，请放大二维码区域重试',
          type: 'warning',
        })
      }
  }
}

// 定时加载生成验证码
function reloadTotp(){
  if (reloadTimeHandle != null){
    clearInterval(reloadTimeHandle)
  }
  totp()
  reloadTimeHandle = setInterval(totp,5000)
}

// 生成totp码
function totp(){
  Totp(data.qdata).then(result=>{
    data.code = result
  }).catch(err => {
    if (reloadTimeHandle != null){
      clearInterval(reloadTimeHandle)
    }
    ElMessage({
      message: err,
      type: 'warning',
    })
  })
}


onMounted(() => {
  // 加载数据文件
  Get().then(result=>{
    data.qdata = result
    if (result != ""){
      reloadTotp()
    } 
  }).catch(err => {
    ElMessage({
        message: '获取数据失败：' + err,
        type: 'warning',
      })
  })
})

</script>

<template>
  <main>

    <div style="padding-top: 10px;">

      <div class="main-header">
        <el-row >
          <el-col :span="6" :offset="6"><div class="grid-content ep-bg-purple" />
            <el-button size="large" type="primary" @click="screenShot" round>扫码添加</el-button>
          </el-col>
          <el-col :span="6"><div class="grid-content ep-bg-purple-light" />
            <el-button size="large" type="primary" @click="showForm" round>手动添加</el-button>
          </el-col>
        </el-row>
      </div>

      <div class="main-body">
          
          <el-row>
            <el-col :span="22" :offset="1">
              <el-card class="box-card">
                <template #header>
                  <div class="card-header">
                    <span>{{ data.username }}</span>
                  </div>
                </template>
                
                <el-row>
                  <el-col :span="10" :offset="6">
                    <div class="code-show">
                      {{ data.code }}
                    </div>
                  </el-col>
                  <el-col :span="2">
                    <el-button type="success" link @click="copy">复制</el-button>
                  </el-col>
                </el-row>

              </el-card>
            </el-col>
          </el-row>

      </div>

    </div>



    <el-dialog v-model="centerDialogVisible" title="请选取二维码" width="80%" center>
      <span>
        <vueCropper style="height:400px;width:100%"
        ref="cropper"
        :img="data.img"
        autoCrop
        centerBox
        outputSize="1"
        ></vueCropper>
      </span>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="centerDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="screenCut">
            确认选取
          </el-button>
        </span>
      </template>
    </el-dialog>


    <el-dialog v-model="addFormVisible" title="手动添加2FA" width="50%" center>
      <el-form 
      :model="form"
      :rules="rules"
      ref="formRef"
      label-width="120px"
      status-icon
      >
        <el-form-item label="账号" prop="account">
          <el-input v-model.trim="form.account" required/>
        </el-form-item>
        <el-form-item label="秘钥" prop="secret">
          <el-input v-model.trim="form.secret" type="password" show-password required/>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addFormVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

  </main>
</template>

<style scoped>

.card-header{
  font-size: 35px;
  color: #A8ABB2;
}
.code-show{
  color: #A8ABB2;
  font-size: 80px;
}

.dialog-footer button:first-child {
  margin-right: 10px;
}

.main-header {
  border-radius: 30px 30px 0 0;
  border-bottom: #E4E7ED 1px solid;
  padding-top: 20px;
  padding-bottom: 20px;
  width: 98%;
  margin:0 auto;
  background-color: #FFFFFF;
}

.main-body{
  border-radius: 0 0 30px 30px;
  padding-top: 20px;
  padding-bottom: 20px;
  width: 98%;
  margin:0 auto;
  background-color: #FFFFFF;
  color: #A8ABB2;
}
</style>
