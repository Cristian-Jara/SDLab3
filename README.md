# SDLab3

### Integrantes

- Cristian Jara 201704563-9

- Sebastian Muoz 201473503-0

- María Riveros 201704585-k

### Instrucciones para funcionamiento

Primero se debe conectar el Broker en la máquina virtual dist86 con ip 10.6.40.227 con el comando "make runBroker", luego se deben conectar todos los servidores en el orden 3,2,1 o 2,3,1 asegurando de dejar al servidor 1 (principal) para el final ya que este se conectará a los otros dos para la consistencia eventual. A continuación se detallará el paso por paso:

1. Conectar el Broker en la máquina virtual dist86 con ip 10.6.40.227 con "make runBroker", en este caso el broker queda solo para cumplir con las restricciones.

2. Para la máquina dist87 con ip 10.6.40.228 conectar el servidor 2 con el comando "make runSv2", en esta misma máquina se debe correr el informante 1 con el comando "make runInf1", de todas formas se recomienda dejar los "clientes" para el final.

3. Para la máquina dist88 con ip 10.6.40.229 conectar el servidor 3 con el comando "make runSv3", de igual manera que antes aquí se conecta el informante 2 con el comando "make runInf2".

4. Para la máquina dist85 con ip 10.6.40.225 se debe conectar el servidor 1 que realiza la función de principal con el comando "make runSv1", en esta misma máquina se conecta a Leia con el comando "make runLeia".

5. Cómo fue mencionado idealmente una vez conectado broker y servidores, deberían conectarse recién los clientes para no generar posibles problemas (como consultar antes de que estos esten conectados).

### Comandos del makefile:

runBroker : Corre el proceso del Broker
runLeia : Corre el proceso de Leia
runInf1 : Corre al informante 1 que corresponde a Ahsoka Tano
runInf2 : Corre al informante 2 que corresponde al almirante Thrawn
runSv1 : Corre el servidor Fulcrum 1 que será el principal
runSv2 : Corre el servidor Fulcrum 2 
runSv3 : Corre el servidor Fulcrum 3 

### Explicación del MERGE y resolución de problemas

Para el merge se realiza por partes tomando la data del servidor 1 (principal) como la inicial desde ahí se va agregando toda la información faltante y actualizando con el mayor número que exista en una ciudad de cierto planeta con la información de los otros servidores. 


De los posibles conflictos, los de AddCity y UpdateNumber fueron ignorados y simplemente se toma el mayor valor de cierta ciudad, para los otros dos eran algo más complejos así que se crearon listas con todos estos conflictos y al final de actualizar la info se arreglaban. Para el updateName se revisaron todos los que tenían el problema y se queda con el último valor actualizado, por otro lado en DeleteCity se asumió que si una ciudad era eliminada en cualquiera de los servidores, sin importar si se actualizaba después u otra cosa se eliminaba.