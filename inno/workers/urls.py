from django.urls import path
from . import views

urlpatterns = [
    path('', views.worker_tree, name='worker_tree'),
    path('department/<int:depId>/', views.worker_tree, name='worker_tree'),
]