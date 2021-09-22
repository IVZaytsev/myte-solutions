from django.contrib import admin
from .models import Position, Worker, Department
# Register your models here.

class WorkerAdmin(admin.ModelAdmin):
    list_display = ['fullname', 'position', 'dateEmploy']

class DepartmentAdmin(admin.ModelAdmin):
    list_display = ['desc', 'get_worker_count']
    filter_horizontal = ('workers',)
    
    def get_worker_count(self, obj):
        return obj.get_worker_count()

    get_worker_count.short_description = 'Кол-во сотрудников'

admin.site.register(Position)
admin.site.register(Worker, WorkerAdmin)
admin.site.register(Department, DepartmentAdmin)