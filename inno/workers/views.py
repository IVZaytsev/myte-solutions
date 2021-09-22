from django.shortcuts import render
from .models import Position, Worker, Department

# Create your views here.
def worker_tree(request, depId: int=None):
    if depId > len(Department.objects.all()): depId = None
    if depId:
        querySet = Department.objects.filter(pk=depId).prefetch_related('workers')
    else:
        querySet = Department.objects.all().prefetch_related('workers')
    return render(request, 'workertree.html', {'departList': querySet})