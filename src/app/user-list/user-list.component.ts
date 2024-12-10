import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

interface User {
  id: number;
  username: string;
  email: string;
  comment: string;
  avatar: string; // 头像路径
}

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports:[FormsModule, CommonModule],
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {
  users: User[] = [];
  newUser: User = { id: 0, username: '', email: '', comment: '', avatar: '' }; 
  apiUrl = 'http://localhost:8080/users'; // 后端 API 地址
  selectedFile: File | null = null; // 存储选中的文件

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.getUsers();
  }

  // 获取用户列表
  getUsers() {
    this.http.get<User[]>(this.apiUrl).subscribe(
      (data) => {
        this.users = data;
      },
      (error) => {
        console.error('获取用户数据出错', error);
      }
    );
  }

  // 监听文件选择
  onFileSelected(event: any) {
    this.selectedFile = event.target.files[0];
  }

  // 添加用户
  addUser() {
    const formData = new FormData();
    formData.append('username', this.newUser.username);
    formData.append('email', this.newUser.email);
    formData.append('comment', this.newUser.comment);

    if (this.selectedFile) {
      formData.append('avatar', this.selectedFile);
    }

    this.http.post<User>(this.apiUrl, formData).subscribe({
      next: (response) => {
        this.users.push(response);
        this.newUser = { id: 0, username: '', email: '', comment: '', avatar: '' }; // 清空表单
        this.selectedFile = null; // 清空文件
        console.log('用户添加成功', response);
      },
      error: (error) => {
        console.error('添加用户出错', error);
      }
    });
  }

  // 删除用户
  deleteUser(username:string) {
    this.http.delete(`${this.apiUrl}/${username}`).subscribe(
      () => {
        this.users = this.users.filter(user => user.username !== username);
      },
      error => console.error('删除用户出错', error)
    );
  }
}
