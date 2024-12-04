import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {
  users: any[] = []; // 存储用户数据

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    this.getUsers(); // 在组件初始化时获取用户数据
  }

  getUsers(): void {
    this.http.get<any[]>('http://localhost:8080/users') // 调用后端 API 获取用户数据
      .subscribe(
        (data) => {
          this.users = data; // 成功获取数据后存储
        },
        (error) => {
          console.error('Error fetching users', error); // 错误处理
        }
      );
  }
}
