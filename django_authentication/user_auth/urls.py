from django.contrib import admin
from django.urls import path, include
from . import views

urlpatterns = [
    path('', views.home_page, name="home"),
    path('signup/', views.SignUp.as_view(), name="signup"),
    path('login/', views.Login.as_view(), name="login"),
    path('password_change', views.PasswordChange.as_view(), name ="password_change"),
    
]